package main

import (
	"fmt"
	"sync"
	"time"

	"context"
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

func (m *Manager) runClient(u musdk.User) {
	logger.Infof("run users %d", u.Port)
	ctx, cancel := context.WithCancel(context.Background())
	client := &Client{
		userID:  u.Id,
		ctx:     ctx,
		cancel:  cancel,
		port:    u.Port,
		dstAddr: fmt.Sprintf("%s:%d", dst, u.Port+Port_Offset),
		lock:    new(sync.Mutex),
	}
	go client.Run()
	logger.Debugf("store client %d to map", u.Port)
	m.clients.Store(u.Port, client)
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

	_, ok := m.clients.Load(u.Port)
	if ok {
		logger.Debugf("%d is running ... skip", u.Port)
		return
	}
	m.runClient(u)
}
