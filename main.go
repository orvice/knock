package main

import (
	"github.com/orvice/kit/log"
)

var (
	logger log.Logger
)

func main() {
	logger = log.NewFileLogger(Log_Path)

	InitEnv()
	InitWebApi()
	m := NewManager()
	m.Daemon()

}
