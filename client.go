package main

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	"context"
	"github.com/catpie/musdk-go"
	"github.com/orvice/utils/log"
)

const (
	MaxErr = 10
)

type Client struct {
	userID int64
	ctx    context.Context
	cancel func()

	dstType string
	dstAddr string
	l       net.Listener
	port    int

	lock *sync.Mutex

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
		select {
		case <-c.ctx.Done():
			return
		default:
			c.handleTcpListener()
		}
	}
}

func (c *Client) handleTcpListener() {
	srcConn, err := c.l.Accept()
	if err != nil {
		return
	}

	d_tcpAddr, _ := net.ResolveTCPAddr("tcp4", c.dstAddr)
	dstConn, err := net.DialTCP("tcp", nil, d_tcpAddr)
	if err != nil {
		logger.Errorf("can't connect " + c.dstAddr)
		srcConn.Close()
		return
	}
	go func() {
		n, err := io.Copy(srcConn, dstConn)
		if err != nil {
			logger.Error("%d src->dst  ", c.port, err)
		}
		go c.AddTraffic(n)
		logger.Infof("%d src->dst Written len: %d", c.port, n)
	}()

	go func() {
		n, err := io.Copy(dstConn, srcConn)
		if err != nil {
			logger.Error("%d dst->src  ", c.port, err)
		}
		go c.AddTraffic(n)
		logger.Infof("%d dst->src Written len: %d", c.port, n)
	}()
}

func (c *Client) AddTraffic(i int64) {
	logger.Debugf("%d add traffic %d", c.port, i)
	apiClient.SaveTrafficLog(musdk.UserTrafficLog{
		UserId: c.userID,
		D:      i,
	})
}
