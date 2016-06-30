package protocol

import ()

type HeartPacket struct {
	ID uint64
}

func (p *HeartPacket) Serialize() []byte {
	var buf []byte
	return buf
}

func ParseHeart(buffer []byte, gatewayid uint64) *HeartPacket {
	return &HeartPacket{
		ID: gatewayid,
	}
}
