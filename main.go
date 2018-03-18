package main

import (
	"fmt"
	"sync"
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
	InitEnv()
	p := startPort
	for p < endPort {
		runClient(p)
		p++
	}
	select {}
}
