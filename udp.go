package main

import (
	"net"
	"time"
)

const BufferSize = 65535

func (c *Client) udp() {
	addr, err := net.ResolveUDPAddr("udp", c.getListenAddr())
	if err != nil {
		logger.Infof("UDP fail to ResolveUDPAddr %s retry in 10s", c.getListenAddr())
		time.Sleep(RetryTime)
		go c.udp()
		return
	}

	logger.Infof("UPD: try to listen on: %s", addr)

	// 监听
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		logger.Infof("UDP fail to listen %s retry in 10s err %v ", c.getListenAddr(), err)
		time.Sleep(RetryTime)
		go c.udp()
		return
	}

	defer conn.Close()

	logger.Infof("UDP try to handle request on: %s", addr)

	for {
		select {
		case <-c.ctx.Done():
			return
		default:
			c.handleConn(conn)
		}
	}
}

func (c *Client) handleConn(conn *net.UDPConn) {
	// 读取请求
	request := make([]byte, BufferSize)
	n, remoteAddr, err := conn.ReadFromUDP(request)
	if err != nil {
		logger.Infof("fail to read from UDP: %v", err)
		return
	}

	logger.Debugf("%d bytes read from %s: \"%s\"", n, remoteAddr.String(), string(request))

	// 转发
	logger.Debugf("try to forward request")
	go c.forward(conn, remoteAddr, request[:n])
}

func (c *Client) forward(conn *net.UDPConn, remoteAddr *net.UDPAddr, request []byte) {
	// 转发请求到目标地址
	logger.Debugf("try to connect target service")
	targetConn, err := net.DialTimeout("udp", c.dstAddr, time.Second*1)

	// io.Copy(conn,targetConn)

	if err != nil {
		logger.Errorf("unable to connect target service with error: %v", err)
		return
	}

	defer targetConn.Close()

	targetConn.SetWriteDeadline(time.Now().Add(time.Second * 1))

	logger.Debugf("UDP try to write request %s to target\n", string(request))
	n, err := targetConn.Write(request)
	if err != nil {
		logger.Errorf("unable to write request to target: %s\n", err.Error())
		return
	}
	logger.Infof("UDP write %d bytes of request to target\n", c.port, n)
	go c.AddTraffic(int64(n))

	data := make([]byte, BufferSize)
	n, err = targetConn.Read(data)
	if err != nil {
		logger.Debugf("UDP unable to read from target: %s\n", err.Error())
		return
	}

	logger.Debugf("udp %d read %d bytes from target: \"%s\"\n", c.port, n, string(data))
	n, err = conn.WriteToUDP(data[:n], remoteAddr)
	if err != nil {
		logger.Errorf("fail to write to remote addr:", remoteAddr, string(request), err)
		return
	}
	c.AddTraffic(int64(n))
	logger.Infof("udp %d %d bytes write to remote addr:", c.port, n, remoteAddr)
}
