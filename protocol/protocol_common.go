package protocol

import (
	"bytes"
	"strings"
)

const (
	ENDFLAG  string = "\r\n"
	SEP      string = ","
	PO       string = "po"
	ST       string = "st"
	XE       string = "xe"
	WA       string = "wa"
	HI       string = "hi"
	RD       string = "rd"
	CG       string = "cg"
	TERM_KEY string = "xexun"
)

var term_keys_list = []int{477, 269, 247, 454, 243, 177, 176, 382, 432, 242, 357, 497, 263, 361, 220, 106, 276, 350, 113, 108, 194, 306, 443, 303, 471, 268, 330, 404, 458, 110, 317, 167, 261, 179, 112, 400, 405, 375, 242, 118, 446, 323, 468, 273, 380, 451, 414, 298, 308, 182, 219, 305, 185, 450, 261, 491, 394, 273, 348, 278, 456, 124, 367, 298, 438, 322, 157, 189, 457, 101, 325, 387, 127, 289, 473, 129, 283, 313, 145, 284, 286, 161, 236, 371, 419, 354, 209, 310, 145, 221}
var CMDIDS = []string{PO, ST, XE, WA, HI, RD, CG}

var (
	Illegal   uint16 = 0
	UnSupport uint16 = 254
	HalfPack  uint16 = 255

	Login     uint16 = 1
	HeartBeat uint16 = 2
	PosUp     uint16 = 3
	WarnUp    uint16 = 4
	Echo      uint16 = 5
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
		return PosUp
	case WA:
		return WarnUp
	case CG:
		return Echo
	case RD:
		return Echo
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
	endindex := strings.Index(cmd, ENDFLAG)
	if endindex == -1 {
		return HalfPack, 0
	} else {
		tmp := cmd[0:endindex]
		cmdindex, cmdid := getCommandIndex(string(tmp))
		if cmdindex != -1 {
			return getCommandID(cmdid), uint16(endindex + 2)
		} else {
			return Illegal, uint16(endindex + 2)
		}
	}

	return HalfPack, 0
}
