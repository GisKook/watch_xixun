package shunt

import (
	"bytes"
	//	"encoding/binary"
	//	"errors"
)

var (
	LeftFlag      byte   = '['
	RightFlag     byte   = ']'
	DasLeftFlag   byte   = '$'
	DasRightFlag1 byte   = '\r'
	DasRightFlag2 byte   = '\n'
	Comma         []byte = []byte{','}
	Asterisk      []byte = []byte{'*'}
	Colon         []byte = []byte{':'}

	protocolID = map[string]uint16{
		"LK": HeartBeat,
		"UD": PosUp,
	}

	protocolIDDas = map[string]uint16{
		"$LOGRT":  Login,
		"$HCHECK": HeartBeat,
	}
)

func parseHeader(buf *bytes.Buffer) uint16 {
	split := bytes.Split(buf.Bytes(), Comma)
	header_strings := bytes.Split(split[0], Asterisk)

	return protocolID[string(header_strings[3])]
}

func CheckProtocol(buffer *bytes.Buffer) (uint16, uint16) {
	bufferlen := buffer.Len()
	if bufferlen == 0 {
		return Illegal, 0
	}
	if buffer.Bytes()[0] != LeftFlag {
		buffer.ReadByte()
		CheckProtocol(buffer)
	}
	for i := 1; i < bufferlen; i++ {
		if buffer.Bytes()[i] == RightFlag {
			return parseHeader(buffer), uint16(i)
		}
	}

	return HalfPack, 0
}

func parseDasCmdID(buf *bytes.Buffer) uint16 {
	split := bytes.Split(buf.Bytes(), Colon)

	return protocolIDDas[string(split[0])]

}

func CheckProtocolDas(buffer *bytes.Buffer) (uint16, uint16) {
	bufferlen := buffer.Len()
	if bufferlen == 0 {
		return Illegal, 0
	}
	if buffer.Bytes()[0] != DasLeftFlag {
		buffer.ReadByte()
		CheckProtocol(buffer)
	}
	for i := 1; i < bufferlen; i++ {
		if buffer.Bytes()[i] == DasRightFlag2 && buffer.Bytes()[i-1] == DasRightFlag1 {
			return parseDasCmdID(buffer), uint16(i)
		}
	}

	return HalfPack, 0

}
