package main

import (
	"github.com/orvice/kit/log"
)

var (
	logger log.Logger
)

func main() {
	InitEnv()
	logger = log.NewFileLogger(Log_Path)

	InitWebApi()
	m := NewManager()
	m.Daemon()

}
