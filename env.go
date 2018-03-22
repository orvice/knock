package main

import (
	"github.com/orvice/utils/env"
	"time"
)

var (
	dst       string
	Mu_Token  string
	Mu_Uri    string
	Mu_NodeID int
	Log_Path  string

	Port_Offset int32

	RetryTime = time.Second * 10
)

func InitEnv() {
	dst = env.Get("DST")
	Mu_Uri = env.Get("MU_URI")
	Mu_Token = env.Get("MU_TOKEN")
	Mu_NodeID = env.GetInt("MU_NODE_ID")

	Port_Offset = int32(env.GetInt("PORT_OFFSET",0))

	Log_Path = env.Get("LOG_PATH","/var/log/knock.log")
}
