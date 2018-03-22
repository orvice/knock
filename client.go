package main

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	"github.com/orvice/utils/log"
	"context"
)

const (
	MaxErr = 10
)

type Client struct {
	ctx context.Context
	cancel func()

	dstType string
	dstAddr string
	l       net.Listener
	port    int32

	lock    *sync.Mutex
	traffic int64

	errCount int64
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
		log.Errorf("listen tcp port %s error %v,will retry in 10s", c.getListenAddr(), err)
		time.Sleep(RetryTime)
		if c.errCount > MaxErr {
			return
		}
		go c.tcp()
		c.errCount++
		return
	}
	logger.Infof("Run client port: %d dstAddr: %s", c.port, c.dstAddr)
	for {
		srcConn, err := c.l.Accept()
		if err != nil {
			continue
		}

		d_tcpAddr, _ := net.ResolveTCPAddr("tcp4", c.dstAddr)
		dstConn, err := net.DialTCP("tcp", nil, d_tcpAddr)
		if err != nil {
			logger.Errorf("can't connect " + c.dstAddr)
			srcConn.Close()
			continue
		}
		go func() {
			n, err := io.Copy(srcConn, dstConn)
			if err != nil {
				logger.Error("src->dst  ", err)
			}
			go c.AddTraffic(n)
			logger.Infof("%d src->dst Written len: %d",c.port, n)
		}()

		go func() {
			n, err := io.Copy(dstConn, srcConn)
			if err != nil {
				logger.Error("dst->src  ", err)
			}
			go c.AddTraffic(n)
			logger.Infof("%d dst->src Written len: %d",c.port, n)
		}()
	}
}

func (c *Client) AddTraffic(i int64) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.traffic += i

	logger.Debugf("Traffic count: %d", c.traffic)
}
