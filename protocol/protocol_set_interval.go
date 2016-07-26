package protocol

import ()

type Set_Interval_Packet struct {
	Encryption   string
	IMEI         string
	SerialNumber string
}

func (p *Set_Interval_Packet) Serialize() []byte {
	var result string
	result += p.Encryption
	result += p.IMEI
	result += SEP
	result += p.SerialNumber
	result += ",123456cmd,"
	result += "ti=30"
	result += ENDFLAG

	return []byte(result)
}

func Parse_Set_Interval(encryption string, imei string, serialnum string) *Set_Interval_Packet {
	return &Set_Interval_Packet{
		Encryption:   encryption,
		IMEI:         imei,
		SerialNumber: serialnum,
	}
}
