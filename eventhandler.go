package watch_xixun

import (
	"github.com/giskook/gotcp"
	"github.com/giskook/watch_xixun/protocol"
	"log"
	"strconv"
	"time"
)

type Callback struct{}

func (this *Callback) OnConnect(c *gotcp.Conn) bool {
	checkinterval := GetConfiguration().GetServerConnCheckInterval()
	readlimit := GetConfiguration().GetServerReadLimit()
	writelimit := GetConfiguration().GetServerWriteLimit()
	config := &ConnConfig{
		ConnCheckInterval: uint16(checkinterval),
		ReadLimit:         uint16(readlimit),
		WriteLimit:        uint16(writelimit),
	}
	conn := NewConn(c, config)

	c.PutExtraData(conn)

	conn.Do()
	NewConns().Add(conn)
	log.Println("<DEBUG> a new connect comes")

	return true
}

func (this *Callback) OnClose(c *gotcp.Conn) {
	conn := c.GetExtraData().(*Conn)
	conn.Close()
	NewConns().Remove(conn.ID)
}

func on_login(c *gotcp.Conn, p *ShaPacket) {
	conn := c.GetExtraData().(*Conn)
	conn.Status = ConnSuccess
	loginPkg := p.Packet.(*protocol.LoginPacket)
	conn.IMEI = loginPkg.IMEI
	conn.ID, _ = strconv.ParseUint(loginPkg.IMEI, 10, 64)
	NewConns().SetID(conn.ID, conn.index)
	c.AsyncWritePacket(p, time.Second)
}

func (this *Callback) OnMessage(c *gotcp.Conn, p gotcp.Packet) bool {
	shaPacket := p.(*ShaPacket)
	log.Println("onmessage")
	switch shaPacket.Type {
	case protocol.Login:
		log.Println("on_login")
		on_login(c, shaPacket)
	case protocol.HeartBeat:
		c.AsyncWritePacket(shaPacket, time.Second)
		//GetServer().GetProducer().Send(GetServer().GetTopic(), p.Serialize())
	}

	return true
}
