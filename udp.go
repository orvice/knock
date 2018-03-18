package main

import (
	"net"
	"time"
)

const BufferSize = 65535

func (c *Client) udp() {
	addr, err := net.ResolveUDPAddr("udp", c.getListenAddr())
	if err != nil {
		panic(err)
	}

	logger.Infof("UPD: try to listen on: %s", addr)

	// 监听
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		panic(err)
	}

	defer conn.Close()

	logger.Infof("UDP try to handle request on: %s", addr)

	for {
		c.handleConn(conn)
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

	logger.Debugf("try to write request %s to target\n", string(request))
	n, err := targetConn.Write(request)
	if err != nil {
		logger.Errorf("unable to write request to target: %s\n", err.Error())
		return
	}
	logger.Debugf("write %d bytes of request to target\n", n)

	data := make([]byte, BufferSize)
	n, err = targetConn.Read(data)
	if err != nil {
		logger.Debugf("unable to read from target: %s\n", err.Error())
		return
	}

	logger.Debugf("read %d bytes from target: \"%s\"\n", n, string(data))
	n, err = conn.WriteToUDP(data[:n], remoteAddr)
	if err != nil {
		logger.Errorf("fail to write to remote addr:", remoteAddr, string(request), err)
		return
	}
	logger.Infof("%d bytes write to remote addr:", n, remoteAddr)
}
