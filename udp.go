package main

import (
	"log"
	"net"
	"os"
	"time"
)

const BufferSize = 65535

var logger = log.New(os.Stdout, "udp => ", log.Ldate|log.Ltime)

func (c *Client) udp() {
	addr, err := net.ResolveUDPAddr("udp", c.getListenAddr())
	if err != nil {
		panic(err)
	}

	logger.Println("try to listen on:", addr)

	// 监听
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		panic(err)
	}

	defer conn.Close()

	logger.Println("try to handle request on:", addr)

	for {
		c.handleConn(conn)
	}
}

func (c *Client) handleConn(conn *net.UDPConn) {
	// 读取请求
	request := make([]byte, BufferSize)
	n, remoteAddr, err := conn.ReadFromUDP(request)
	if err != nil {
		logger.Println("fail to read from UDP:", err)
		return
	}

	logger.Printf("%d bytes read from %s: \"%s\"", n, remoteAddr.String(), string(request))

	// 转发
	logger.Println("try to forward request")
	go c.forward(conn, remoteAddr, request[:n])
}

func (c *Client) forward(conn *net.UDPConn, remoteAddr *net.UDPAddr, request []byte) {
	// 转发请求到目标地址
	logger.Println("try to connect target service")
	targetConn, err := net.DialTimeout("udp", c.dstAddr, time.Second*1)

	// io.Copy(conn,targetConn)

	if err != nil {
		logger.Println("unable to connect target service with error:", err)
		return
	}

	defer targetConn.Close()

	targetConn.SetWriteDeadline(time.Now().Add(time.Second * 1))

	logger.Printf("try to write request %s to target\n", string(request))
	n, err := targetConn.Write(request)
	if err != nil {
		logger.Printf("unable to write request to target: %s\n", err.Error())
		return
	}
	logger.Printf("write %d bytes of request to target\n", n)

	data := make([]byte, BufferSize)
	n, err = targetConn.Read(data)
	if err != nil {
		logger.Printf("unable to read from target: %s\n", err.Error())
		return
	}

	logger.Printf("read %d bytes from target: \"%s\"\n", n, string(data))
	n, err = conn.WriteToUDP(data[:n], remoteAddr)
	if err != nil {
		logger.Println("fail to write to remote addr:", remoteAddr, string(request), err)
		return
	}
	logger.Println(n, "bytes write to remote addr:", remoteAddr)
}
