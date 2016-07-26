package protocol

import ()

type Action_Packet struct {
	Encryption   string
	IMEI         string
	SerialNumber string
	Action       string
}

func (p *Action_Packet) Serialize() []byte {
	var result string
	result += p.Encryption
	result += p.IMEI
	result += SEP
	result += p.SerialNumber
	result += ",123456cmd,"
	result += p.Action
	result += ENDFLAG

	return []byte(result)
}

func Parse_Action(encryption string, imei string, serialnum string, action string) *Action_Packet {
	return &Action_Packet{
		Encryption:   encryption,
		IMEI:         imei,
		SerialNumber: serialnum,
		Action:       action,
	}
}
