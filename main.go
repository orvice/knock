package main

import (
	"fmt"
	"github.com/orvice/kit/log"
	"sync"
)

var (
	logger log.Logger
)

func runClient(port int32) {
	client := &Client{
		port:    port,
		dstAddr: fmt.Sprintf("%s:%d", dst, port),
		lock:    new(sync.Mutex),
	}
	go client.Run()
}

func main() {
	logger = log.NewDefaultLogger()
	InitEnv()
	p := startPort
	for p < endPort {
		runClient(p)
		p++
	}
	select {}
}
