package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/catpie/musdk-go"
)

type Manager struct {
	clients *sync.Map
}

func NewManager() *Manager {
	return &Manager{
		clients: new(sync.Map),
	}
}

func (m *Manager) runClient(port int32) {
	logger.Infof("run users %d", port)
	client := &Client{
		port:    port,
		dstAddr: fmt.Sprintf("%s:%d", dst, port),
		lock:    new(sync.Mutex),
	}
	go client.Run()
	logger.Infof("store client %d to map", port)
	m.clients.Store(port, client)
}

func (m *Manager) Daemon() {
	for {
		users, err := apiClient.GetUsers()
		if err != nil {
			logger.Error(err)
			time.Sleep(time.Minute)
		}
		for _, u := range users {
			go m.checkUser(u)
		}
		time.Sleep(time.Second * 300)
	}
}

func (m *Manager) checkUser(u musdk.User) {

	if !u.IsEnable() {
		return
	}

	_, ok := m.clients.Load(int32(u.Port))
	if ok {
		logger.Infof("%d is running ... skip", u.Port)
		return
	}
	m.runClient(int32(u.Port))
}
