package watch_xixun

import (
	"bytes"
	"github.com/giskook/gotcp"
	"log"
	"time"
)

var ConnSuccess uint8 = 0
var ConnUnauth uint8 = 1

type ConnConfig struct {
	ConnCheckInterval uint16
	ReadLimit         uint16
	WriteLimit        uint16
	NsqChanLimit      uint16
}

type Conn struct {
	conn          *gotcp.Conn
	config        *ConnConfig
	recieveBuffer *bytes.Buffer
	ticker        *time.Ticker
	readflag      int64
	writeflag     int64
	closeChan     chan bool
	index         uint32
	ID            uint64
	IMEI          string
	Status        uint8
	ReadMore      bool
}

func NewConn(conn *gotcp.Conn, config *ConnConfig) *Conn {
	return &Conn{
		conn:          conn,
		recieveBuffer: bytes.NewBuffer([]byte{}),
		config:        config,
		readflag:      time.Now().Unix(),
		writeflag:     time.Now().Unix(),
		ticker:        time.NewTicker(time.Duration(config.ConnCheckInterval) * time.Second),
		closeChan:     make(chan bool),
		index:         0,
		Status:        ConnUnauth,
		ReadMore:      true,
	}
}

func (c *Conn) Close() {
	c.closeChan <- true
	c.ticker.Stop()
	c.recieveBuffer.Reset()
	close(c.closeChan)
}

func (c *Conn) GetBuffer() *bytes.Buffer {
	return c.recieveBuffer
}

func (c *Conn) writeToclientLoop() {
	defer func() {
		c.conn.Close()
	}()

	for {
		select {
		case <-c.closeChan:
			return
		}
	}
}

func (c *Conn) SendToClient(p gotcp.Packet) {
	c.conn.AsyncWritePacket(p, time.Second)
}

func (c *Conn) UpdateReadflag() {
	c.readflag = time.Now().Unix()
}

func (c *Conn) UpdateWriteflag() {
	c.writeflag = time.Now().Unix()
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
			if now-c.readflag > int64(c.config.ReadLimit) {
				log.Println("read linmit")
				return
			}
			//			if c.Status == ConnUnauth {
			//				log.Printf("unauth's gateway gatewayid %x\n", c.ID)
			//				return
			//			}
		case <-c.closeChan:
			return
		}
	}
}

func (c *Conn) Do() {
	go c.checkHeart()
	go c.writeToclientLoop()
}
