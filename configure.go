package watch_xixun

import (
	"encoding/json"
	"os"
)

type ServerConfiguration struct {
	ReadLimit         int64
	WriteLimit        int64
	ConnTimeout       uint8
	ConnCheckInterval uint8
	ServerStatistics  uint32
	BindPort          string
}

type DatabaseConfiguration struct {
	Host   string
	Port   string
	User   string
	Passwd string
	Dbname string
}

type Configuration struct {
	ServerConfig *ServerConfiguration
	DBConfig     *DatabaseConfiguration
}

func ReadConfig(confpath string) (*Configuration, error) {
	file, _ := os.Open(confpath)
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err := decoder.Decode(&configuration)

	return &configuration, err
}

func (conf *Configuration) GetServerReadLimit() int64 {
	return conf.ServerConfig.ReadLimit
}

func (conf *Configuration) GetServerWriteLimit() int64 {
	return conf.ServerConfig.WriteLimit
}

func (conf *Configuration) GetServerConnCheckInterval() uint8 {
	return conf.ServerConfig.ConnCheckInterval
}

func (conf *Configuration) GetServerStatistics() uint32 {
	return conf.ServerConfig.ServerStatistics
}

var Config *Configuration

func SetConfiguration(config *Configuration) {
	Config = config
}

func GetConfiguration() *Configuration {
	return Config
}
