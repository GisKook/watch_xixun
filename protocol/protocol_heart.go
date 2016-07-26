package protocol

import ()

type HeartPacket struct {
	Content []byte
}

func (p *HeartPacket) Serialize() []byte {
	return p.Content
}

func ParseHeart(buffer []byte) *HeartPacket {
	return &HeartPacket{
		Content: buffer,
	}
}
