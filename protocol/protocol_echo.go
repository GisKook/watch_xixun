package protocol

import ()

type EchoPacket struct {
	Content []byte
}

func (p *EchoPacket) Serialize() []byte {
	return p.Content
}

func ParseEcho(buffer []byte) *EchoPacket {
	return &EchoPacket{
		Content: buffer,
	}
}
