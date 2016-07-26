package protocol

import ()

type Locate_Packet struct {
	Encryption   string
	IMEI         string
	SerialNumber string
}

func (p *Locate_Packet) Serialize() []byte {
	var result string
	result += p.Encryption
	result += p.IMEI
	result += SEP
	result += p.SerialNumber
	result += ",123456cmd,"
	result += "tk=0"
	result += ENDFLAG

	return []byte(result)
}

func Parse_Locate(encryption string, imei string, serialnum string) *Locate_Packet {
	return &Locate_Packet{
		Encryption:   encryption,
		IMEI:         imei,
		SerialNumber: serialnum,
	}
}
