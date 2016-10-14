package protocol

import (
	"log"
	"strconv"
)

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
	Wifi         string
	WifiCount    int
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

func ParseWifi(wifis []string) string {
	log.Println(wifis)
	item_count := len(wifis)
	var ret string
	var j int = 0
	for i := 0; i < item_count; i++ {
		ret += wifis[j] + ","
		ret += wifis[j+1] + ","
		ret += "TP_LINK|"
		j += 2
	}

	ret = ret[0 : len(ret)-1]

	return ret
}

func ParsePosUp(buffer []byte) *PosUpPacket {
	log.Println("parseposupdata")
	encryption, values := ParseCommon(buffer)

	lat := values[5][1:]
	long := values[6][1:]
	wifi_count := values[12]
	log.Println(wifi_count)
	count, _ := strconv.Atoi(wifi_count)
	var wifis string = ""
	if count > 1 {
		wifis = ParseWifi(values[13 : 13+count*2])
	}

	return &PosUpPacket{
		Encryption:   encryption,
		IMEI:         values[1],
		SerialNumber: values[2],
		LocationTime: values[3],
		Latitude:     lat,
		Longitude:    long,
		GPSFlag:      values[len(values)-7],
		Wifi:         wifis,
		WifiCount:    count,
	}
}
