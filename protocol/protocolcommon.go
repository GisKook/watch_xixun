package protocol

import (
	"bytes"
	"log"
	"strings"
)

const (
	ENDFLAG string = "\r\n"
	SEP     string = ","
	PO      string = "po"
	ST      string = "st"
	XE      string = "xe"
	WA      string = "wa"
	HI      string = "hi"
	RD      string = "rd"
	CG      string = "cg"
)

var CMDIDS = []string{PO, ST, XE, WA, HI, RD, CG}

var (
	Illegal   uint16 = 0
	UnSupport uint16 = 254
	HalfPack  uint16 = 255

	Login     uint16 = 1
	HeartBeat uint16 = 2
	PosUP     uint16 = 3
	WarnUP    uint16 = 4
)

func ParseCommon(buffer []byte) (string, []string) {
	tmp := string(buffer[16:])
	value := strings.Split(tmp, SEP)

	return string(buffer[0:16]), value
}

func getCommandIndex(cmd string) (int, string) {
	for i := 0; i < len(CMDIDS); i++ {
		index := strings.Index(cmd, CMDIDS[i])
		if index != -1 {
			return index, CMDIDS[i]
		}
	}

	return -1, ""
}

func getCommandID(cmdid string) uint16 {
	switch cmdid {
	case PO:
		return Login
	case HI:
		return HeartBeat
	case XE:
		return PosUP
	case WA:
		return WarnUP
	default:
		return UnSupport
	}
}

func CheckProtocol(buffer *bytes.Buffer) (uint16, uint16) {
	//log.Printf("check protocol %x\n", buffer.Bytes())
	bufferlen := buffer.Len()
	if bufferlen == 0 {
		return Illegal, 0
	}

	cmd := string(buffer.Bytes()[:bufferlen])
	log.Println(cmd)
	endindex := strings.Index(cmd, ENDFLAG)
	log.Println(endindex)
	if endindex == -1 {
		return HalfPack, 0
	} else {
		tmp := cmd[0:endindex]
		log.Println(tmp)
		cmdindex, cmdid := getCommandIndex(string(tmp))
		if cmdindex != -1 {
			return getCommandID(cmdid), uint16(endindex + 2)
		} else {
			return Illegal, uint16(endindex + 2)
		}
	}

	return HalfPack, 0
}
