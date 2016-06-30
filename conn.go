package shunt

import (
	"bytes"
	"fmt"
	"github.com/giskook/gotcp"
	"github.com/giskook/shunt/protocol"
	"log"
	"net"
	"time"
)

var ConnUnauth uint8 = 0
var ConnSuccess uint8 = 1

type ConnConfig struct {
	HeartBeat  uint8
	ReadLimit  int64
	WriteLimit int64
}

type Conn struct {
	conn         *gotcp.Conn
	config       *ConnConfig
	RecvBuffer   *bytes.Buffer
	ticker       *time.Ticker
	readflag     int64
	writeflag    int64
	closeChan    chan bool
	index        uint32
	IMEI         uint64
	Batt         string
	Status       uint8
	ReadMore     bool
	SetLocale    bool
	Manufacturer string
	dasticker    *time.Ticker

	dasconn       *net.TCPConn
	dasCmdChan    chan gotcp.Packet
	RecvBufferDas *bytes.Buffer
	ReadMoreDas   bool
}

func NewConn(conn *gotcp.Conn, config *ConnConfig) *Conn {
	log.Println(GetConfiguration().GetDasHost())
	tcpaddr, _ := net.ResolveTCPAddr("tcp", GetConfiguration().GetDasHost())
	dasconn, err := net.DialTCP("tcp", nil, tcpaddr)
	if err != nil {
		log.Printf("conn to das fail %s\n", err.Error())
		return nil
	}
	return &Conn{
		conn:       conn,
		RecvBuffer: bytes.NewBuffer([]byte{}),
		config:     config,
		readflag:   time.Now().Unix(),
		writeflag:  time.Now().Unix(),
		ticker:     time.NewTicker(time.Duration(config.HeartBeat) * time.Second),
		dasticker:  time.NewTicker(time.Duration(120 * time.Second)),
		closeChan:  make(chan bool),
		index:      0,
		Status:     ConnUnauth,
		ReadMore:   false,
		SetLocale:  false,

		dasconn:       dasconn,
		dasCmdChan:    make(chan gotcp.Packet, 64),
		RecvBufferDas: bytes.NewBuffer([]byte{}),
		ReadMoreDas:   false,
	}
}

func (c *Conn) Close() {
	c.closeChan <- true
	c.ticker.Stop()
	c.dasticker.Stop()
	c.RecvBuffer.Reset()
	c.RecvBufferDas.Reset()
	c.dasconn.Close()
	c.dasconn = nil
	close(c.closeChan)
	close(c.dasCmdChan)
}

func (c *Conn) sendToDas() {
	defer func() {
		c.conn.Close()
		if c.dasconn != nil {
			c.dasconn.Close()
		}
	}()

	for {
		select {
		case p := <-c.dasCmdChan:
			if p != nil {
				c.dasconn.Write(p.Serialize())
			}
		case <-c.closeChan:
			return
		}
	}
}

func (c *Conn) WriteToDas(p gotcp.Packet) {
	c.dasCmdChan <- p
}

func (c *Conn) UpdateReadflag() {
	c.readflag = time.Now().Unix()
}

func (c *Conn) UpdateWriteflag() {
	c.writeflag = time.Now().Unix()
}

func (c *Conn) SetStatus(status uint8) {
	c.Status = status
}

func (c *Conn) checkHeart() {
	defer func() {
		c.conn.Close()
	}()

	var now int64
	for {
		select {
		case <-c.ticker.C:
			now = time.Now().Unix()
			if now-c.readflag > c.config.ReadLimit {
				log.Println("read linmit")
				return
			}
			if now-c.writeflag > c.config.WriteLimit {
				log.Println("write limit")
				return
			}
			if c.Status == ConnUnauth {
				log.Printf("unauth's gateway gatewayid %d\n", c.IMEI)
				return
			}
		case <-c.closeChan:
			return
		}
	}
}

func (c *Conn) sendDasHeart() {
	defer func() {
		c.conn.Close()
	}()
	for {
		select {
		case <-c.dasticker.C:
			heartpkg := protocol.NewDasHeartPacket(fmt.Sprint(c.IMEI), c.Batt)
			c.WriteToDas(heartpkg)
			log.Println("<OUT DAS> " + string(heartpkg.Serialize()))
		case <-c.closeChan:
			return
		}
	}
}

func (c *Conn) recvdas() {
	for {
		if c.dasconn != nil {
			if c.ReadMoreDas {
				buffer := make([]byte, 1024)
				readLength, err := c.dasconn.Read(buffer)
				if err != nil {
					return
				}

				log.Println("<IN DAS>  " + string(buffer))
				c.RecvBufferDas.Write(buffer[0:readLength])
			}
			cmdid, pkglen := CheckProtocolDas(c.RecvBufferDas)
			pkgbyte := make([]byte, pkglen)
			c.RecvBufferDas.Read(pkgbyte)
			switch cmdid {
			case Login:
				loginRt := protocol.ParseDasLoginRt(pkgbyte)
				if loginRt.Result {
					c.Status = ConnSuccess
					heartpkg := protocol.NewDasHeartPacket(fmt.Sprint(c.IMEI), c.Batt)
					c.WriteToDas(heartpkg)
					log.Println("<OUT DAS> " + string(heartpkg.Serialize()))
				}
				c.ReadMoreDas = false
			case HeartBeat:
				c.ReadMoreDas = false
			case Illegal:
				c.ReadMoreDas = true
			case HalfPack:
				c.ReadMoreDas = true
			}
		} else {
			return
		}
	}
}

func (c *Conn) Do() {
	//go c.checkHeart()
	go c.recvdas()
	go c.sendToDas()
	go c.sendDasHeart()
	//go c.writeToclientLoop()
}

type Callback struct{}

func (this *Callback) OnConnect(c *gotcp.Conn) bool {
	log.Println("new conn ")
	heartbeat := GetConfiguration().GetServerConnCheckInterval()
	readlimit := GetConfiguration().GetServerReadLimit()
	writelimit := GetConfiguration().GetServerWriteLimit()
	config := &ConnConfig{
		HeartBeat:  uint8(heartbeat),
		ReadLimit:  int64(readlimit),
		WriteLimit: int64(writelimit),
	}
	//log.Println(heartbeat,readlimit,writelimit)
	conn := NewConn(c, config)
	if conn != nil {
		c.PutExtraData(conn)
		conn.Do()

		return true
	}

	return false
}

func (this *Callback) OnClose(c *gotcp.Conn) {
	log.Println("closeconn")
	conn := c.GetExtraData().(*Conn)
	conn.Close()
	NewConns().Remove(conn.IMEI)
}

func (this *Callback) OnMessage(c *gotcp.Conn, p gotcp.Packet) bool {
	trackerPacket := p.(*TrackerPacket)
	switch trackerPacket.Type {
	case Login:
		c.AsyncWritePacket(trackerPacket, time.Second)
	case HeartBeat:
		c.AsyncWritePacket(trackerPacket, time.Second)
	}

	return true
}
