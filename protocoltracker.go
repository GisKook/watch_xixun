package shunt

import (
	"github.com/giskook/gotcp"
	"github.com/giskook/shunt/protocol"

	"log"
)

var (
	Illegal   uint16 = 0
	HalfPack  uint16 = 255
	UnSupport uint16 = 254

	Login     uint16 = 1
	HeartBeat uint16 = 2
	PosUp     uint16 = 3
	SetLocale uint16 = 4
)

type TrackerPacket struct {
	Type   uint16
	Packet gotcp.Packet
}

func (this *TrackerPacket) Serialize() []byte {
	switch this.Type {
	case Login:
		return this.Packet.(*protocol.LoginPacket).Serialize()
	case HeartBeat:
		return this.Packet.(*protocol.HeartPacket).Serialize()
	case SetLocale:
		return this.Packet.(*protocol.LocalePacket).Serialize()
	}

	return nil
}

func NewTrackerPacket(Type uint16, Packet gotcp.Packet) *TrackerPacket {
	return &TrackerPacket{
		Type:   Type,
		Packet: Packet,
	}
}

type TrackerProtocol struct {
}

func (this *TrackerProtocol) ReadPacket(c *gotcp.Conn) (gotcp.Packet, error) {
	smconn := c.GetExtraData().(*Conn)
	smconn.UpdateReadflag()

	buffer := smconn.RecvBuffer
	conn := c.GetRawConn()
	for {
		if smconn.ReadMore {
			data := make([]byte, 2048)
			readLengh, err := conn.Read(data)
			log.Println("<IN>      " + string(data))
			if err != nil {
				return nil, err
			}

			if readLengh == 0 {
				return nil, gotcp.ErrConnClosing
			}
			buffer.Write(data[0:readLengh])
		}
		cmdid, pkglen := CheckProtocol(buffer)
		if cmdid == HeartBeat && smconn.Status == ConnUnauth {
			cmdid = Login
		}
		pkgbyte := make([]byte, pkglen)
		buffer.Read(pkgbyte)
		switch cmdid {
		case Login:
			pkg, daspkg, imei, batt, manufacturer := protocol.ParseLogin(pkgbyte)
			smconn.IMEI = imei
			smconn.Batt = batt
			smconn.Manufacturer = manufacturer
			smconn.WriteToDas(daspkg)

			smconn.ReadMore = false
			log.Println("<OUT DAS> " + string(daspkg.Serialize()))
			log.Println("<OUT>     " + string(pkg.Serialize()))
			return NewTrackerPacket(Login, pkg), nil
		case HeartBeat:
			pkg, daspkg, batt := protocol.ParseHeart(pkgbyte)
			smconn.WriteToDas(daspkg)
			smconn.ReadMore = false
			smconn.Batt = batt
			log.Println("<OUT DAS> " + string(daspkg.Serialize()))
			log.Println("<OUT>     " + string(pkg.Serialize()))
			return NewTrackerPacket(HeartBeat, pkg), nil
		case PosUp:
			daspkg := protocol.ParsePosUp(pkgbyte)
			smconn.WriteToDas(daspkg)
			log.Println("<OUT DAS> " + string(daspkg.Serialize()))
			smconn.ReadMore = false
		case Illegal:
			smconn.ReadMore = true
		case HalfPack:
			smconn.ReadMore = true
		}
	}

}
