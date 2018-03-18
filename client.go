package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"net"
	"sync"
)

type Client struct {
	dstType string
	dstAddr string
	l       net.Listener
	port    int32

	lock    *sync.Mutex
	traffic int64
}

func (c *Client) getListenAddr() string {
	return fmt.Sprintf("0.0.0.0:%d", c.port)
}

func (c *Client) Run() {
	go c.tcp()
	go c.udp()
}

func (c *Client) tcp() {
	var err error
	c.l, err = net.Listen("tcp", c.getListenAddr())
	if err != nil {
		fmt.Println(err, err.Error())
	}
	fmt.Printf("Run client port: %d dstAddr: %s", c.port, c.dstAddr)
	for {
		srcConn, err := c.l.Accept()
		if err != nil {
			continue
		}

		d_tcpAddr, _ := net.ResolveTCPAddr("tcp4", c.dstAddr)
		dstConn, err := net.DialTCP("tcp", nil, d_tcpAddr)
		if err != nil {
			fmt.Println(err)
			fmt.Print("can't connect " + c.dstAddr)
			srcConn.Close()
			continue
		}
		go func() {
			n, err := io.Copy(srcConn, dstConn)
			if err != nil {
				log.Error("src->dst  ", err)
			}
			go c.AddTraffic(n)
			log.Infof("src->dst Written len: %d", n)
		}()

		go func() {
			n, err := io.Copy(dstConn, srcConn)
			if err != nil {
				log.Error("dst->src  ", err)
			}
			go c.AddTraffic(n)
			log.Infof("dst->src Written len: %d", n)
		}()
	}
}

func (c *Client) AddTraffic(i int64) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.traffic += i

	log.Infof("Traffic count: %d", c.traffic)
}
