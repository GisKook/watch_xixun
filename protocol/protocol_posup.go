package protocol

import ()

const (
	FEEDBACK_CMDID string = "ac"
)

type PosUpPacket struct {
	Encryption   string
	IMEI         string
	SerialNumber string
	LocationTime string
	Longitude    string
	Latitude     string
	GPSFlag      string
}

func (p *PosUpPacket) Serialize() []byte {
	var result string
	result += p.Encryption
	result += FEEDBACK_CMDID
	result += SEP
	result += p.IMEI
	result += SEP
	result += p.SerialNumber
	result += SEP
	if p.GPSFlag != "" {
		result += "1,"
	} else {
		result += "0,"
	}
	result += p.LocationTime
	result += SEP
	result += ENDFLAG

	return []byte(result)
}

func ParsePosUp(buffer []byte) *PosUpPacket {
	encryption, values := ParseCommon(buffer)

	lat := values[5][1:]
	long := values[6][1:]

	return &PosUpPacket{
		Encryption:   encryption,
		IMEI:         values[1],
		SerialNumber: values[2],
		LocationTime: values[3],
		Latitude:     lat,
		Longitude:    long,
		GPSFlag:      values[len(values)-7],
	}
}
