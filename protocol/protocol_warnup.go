package protocol

type WarnUpPacket struct {
	Encryption   string
	IMEI         string
	SerialNumber string
	WarnStyle    string
}

func (p *WarnUpPacket) Serialize() []byte {
	var result string
	result += p.Encryption
	result += WA
	result += SEP
	result += p.IMEI
	result += SEP
	result += p.SerialNumber
	result += SEP
	result += p.WarnStyle
	result += SEP
	result += ENDFLAG

	return []byte(result)
}

func ParseWarnUp(buffer []byte) *WarnUpPacket {
	encryption, values := ParseCommon(buffer)
	return &WarnUpPacket{
		Encryption:   encryption,
		IMEI:         values[1],
		SerialNumber: values[2],
		WarnStyle:    values[3],
	}
}
