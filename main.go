package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/catpie/musdk-go"
	"github.com/orvice/kit/log"
)

var (
	logger log.Logger
	ports  *sync.Map
)

func runClient(port int32) {
	logger.Infof("run users %d", port)
	client := &Client{
		port:    port,
		dstAddr: fmt.Sprintf("%s:%d", dst, port),
		lock:    new(sync.Mutex),
	}
	go client.Run()
	ports.Store(port, client)
}

func main() {
	logger = log.NewFileLogger(Log_Path)
	ports = new(sync.Map)

	InitEnv()
	InitWebApi()
	Daemon()

}

func Daemon() {
	for {
		users, err := apiClient.GetUsers()
		if err != nil {
			logger.Error(err)
			time.Sleep(time.Minute)
		}
		for _, u := range users {
			checkUser(u)
		}
		time.Sleep(time.Second * 300)
	}
}

func checkUser(u musdk.User) {

	if !u.IsEnable() {
		return
	}

	_, ok := ports.Load(u.Port)
	if ok {
		return
	}
	runClient(int32(u.Port))
}
