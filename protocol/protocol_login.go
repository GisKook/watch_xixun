package protocol

import (
	"strconv"
	"time"
)

type LoginPacket struct {
	Encryption   string
	IMEI         string
	SerialNumber string
}

//var defaultChs = []byte{'2', '3', '4', '5', '6', '7', '8', '9', 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'm', 'n', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z'}
var sys = "23m4iy8hgvjteux5d7p9nzf6qabwcksr"

func generateValiCode(imei string) string {
	imeilen := len(imei)
	start := imei[6]
	mid := imei[7]
	end := imei[8]
	var code string
	for i := 0; i < 6; i++ {
		val := int((imei[i] + imei[imeilen-1-i] + start + end)) * int(mid) * (i + 1)
		tmp, _ := strconv.Atoi(imei[10:])
		result := 31*tmp + val
		result = result % len(sys)
		code += string(sys[result])
	}

	return code
}

func (p *LoginPacket) Serialize() []byte {
	var result string
	result += p.Encryption
	result += p.IMEI
	result += SEP
	result += p.SerialNumber
	result += ",123456cmd,"
	result += "te="
	t := time.Now()
	timelogin := t.Format("060102150405")
	result += timelogin
	result += SEP
	result += generateValiCode(p.IMEI)
	result += ENDFLAG

	return []byte(result)
}

func ParseLogin(buffer []byte) *LoginPacket {
	encryption, values := ParseCommon(buffer)
	return &LoginPacket{
		Encryption:   encryption,
		IMEI:         values[1],
		SerialNumber: values[2],
	}
}
