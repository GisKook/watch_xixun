package main

import (
	"fmt"
	"github.com/giskook/gotcp"
	"github.com/giskook/watch_xixun"
	"log"
	"net"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	// read configuration
	configuration, err := watch_xixun.ReadConfig("./conf.json")
	watch_xixun.SetConfiguration(configuration)

	checkError(err)
	// creates a tcp listener
	tcpAddr, err := net.ResolveTCPAddr("tcp4", ":"+configuration.ServerConfig.BindPort)
	checkError(err)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)

	// creates a tcp server
	config := &gotcp.Config{
		PacketSendChanLimit:    20,
		PacketReceiveChanLimit: 20,
	}
	srv := gotcp.NewServer(config, &watch_xixun.Callback{}, &watch_xixun.ShaProtocol{})

	// create db server
	dbsrv, _ := watch_xixun.NewDbServer(configuration.DBConfig)
	// create watch_xixun server
	watch_xixunserverconfig := &watch_xixun.ServerConfig{
		Listener:      listener,
		AcceptTimeout: time.Duration(configuration.ServerConfig.ConnTimeout) * time.Second,
	}
	watch_xixunserver := watch_xixun.NewServer(srv, watch_xixunserverconfig, dbsrv)
	watch_xixun.SetServer(watch_xixunserver)
	//create watch_data
	watchdata := watch_xixun.NewWatchData(dbsrv)
	go watchdata.Do()

	// starts service
	fmt.Println("listening:", listener.Addr())
	watch_xixunserver.Start()

	// catchs system signal
	chSig := make(chan os.Signal)
	signal.Notify(chSig, syscall.SIGINT, syscall.SIGTERM)
	fmt.Println("Signal: ", <-chSig)

	// stops service
	watch_xixunserver.Stop()
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
