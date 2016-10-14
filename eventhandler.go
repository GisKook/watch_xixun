package watch_xixun

import (
	"fmt"
	"github.com/giskook/gotcp"
	"github.com/giskook/watch_xixun/protocol"
	"log"
	"strconv"
	//	"strings"
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

	return true
}

func (this *Callback) OnClose(c *gotcp.Conn) {
	conn := c.GetExtraData().(*Conn)
	NewConns().Remove(conn)
	conn.Close()
}

func on_login(c *gotcp.Conn, p *ShaPacket) {
	conn := c.GetExtraData().(*Conn)
	conn.Status = ConnSuccess
	loginPkg := p.Packet.(*protocol.LoginPacket)
	conn.IMEI = loginPkg.IMEI
	conn.ID, _ = strconv.ParseUint(loginPkg.IMEI, 10, 64)
	NewConns().SetID(conn.ID, conn.index)
	c.AsyncWritePacket(p, time.Second)
	time.AfterFunc(1*time.Second, func() {
		set_interval_pkg := protocol.Parse_Set_Interval(loginPkg.Encryption, loginPkg.IMEI, loginPkg.SerialNumber)
		c.AsyncWritePacket(set_interval_pkg, time.Second)
	})
}

func on_posup(c *gotcp.Conn, p *ShaPacket) {
	log.Println("posupin")
	posup_pkg := p.Packet.(*protocol.PosUpPacket)
	if posup_pkg.WifiCount > 2 {
		posup_pkg.LocationTime = time.Now().Format("060102150405")
		sql := fmt.Sprintf("INSERT INTO t_posup_log(id, imme,location_time,datainfo,accesstype) VALUES ('%s',to_timestamp('%s','YYMMDDhh24miss'),'%s','%s')",
			posup_pkg.IMEI, posup_pkg.LocationTime, posup_pkg.Wifi, 1)
		log.Println("heihei", sql)
		GetServer().Dbsrv.Insert(sql)
		c.AsyncWritePacket(p, time.Second)
	} else if posup_pkg.GPSFlag != "" {
		c.AsyncWritePacket(p, time.Second)
		log.Println("long", posup_pkg.Longitude)
		if posup_pkg.Longitude != "" {
			log.Println("-----tag 2")
			posup_pkg.LocationTime = time.Now().Format("060102150405")
			sql := fmt.Sprintf("INSERT INTO t_posup_log(imme,location_time,glat,glong) VALUES ('%s',to_timestamp('%s','YYMMDDhh24miss'),'%s','%s')",
				posup_pkg.IMEI, posup_pkg.LocationTime, posup_pkg.Latitude, posup_pkg.Longitude)
			log.Println(sql)
			GetServer().Dbsrv.Insert(sql)
		} else if posup_pkg.WifiCount > 2 {
			log.Println("-----tag 3")
			posup_pkg.LocationTime = time.Now().Format("060102150405")
			sql := fmt.Sprintf("INSERT INTO t_posup_log(imme,location_time,datainfo,accesstype) VALUES ('%s',to_timestamp('%s','YYMMDDhh24miss'),'%s','%s')",
				posup_pkg.IMEI, posup_pkg.LocationTime, posup_pkg.Wifi, 1)
			log.Println("heihei", sql)
			GetServer().Dbsrv.Insert(sql)
		}
	}
}

func on_warnup(c *gotcp.Conn, p *ShaPacket) {
	c.AsyncWritePacket(p, time.Second)

	warnup_pkg := p.Packet.(*protocol.WarnUpPacket)
	sql := fmt.Sprintf("INSERT INTO t_warnup_log(imme,warnstyle,warn_time) VALUES ('%s','%s',to_timestamp('%s','YYMMDDhh24miss'))",
		warnup_pkg.IMEI, warnup_pkg.WarnStyle, time.Now().Format("060102150405"))
	log.Println(sql)
	GetServer().Dbsrv.Insert(sql)

	time.AfterFunc(1*time.Second, func() {
		locate_pkg := protocol.Parse_Locate(warnup_pkg.Encryption, warnup_pkg.IMEI, warnup_pkg.SerialNumber)
		c.AsyncWritePacket(locate_pkg, time.Second)
	})

}

func (this *Callback) OnMessage(c *gotcp.Conn, p gotcp.Packet) bool {
	shaPacket := p.(*ShaPacket)
	log.Println("onmessage packettype", shaPacket.Type)
	switch shaPacket.Type {
	case protocol.Login:
		on_login(c, shaPacket)
	case protocol.HeartBeat:
		c.AsyncWritePacket(shaPacket, time.Second)
	case protocol.PosUp:
		on_posup(c, shaPacket)
	case protocol.Echo:
		c.AsyncWritePacket(shaPacket, time.Second)
	case protocol.WarnUp:
		on_warnup(c, shaPacket)
	}

	return true
}
