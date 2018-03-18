package main

import "github.com/orvice/utils/env"

var (
	dst       string
	startPort int32
	endPort   int32
)

func InitEnv() {
	dst = env.Get("DST")
	startPort = int32(env.GetInt("START_PORT"))
	endPort = int32(env.GetInt("END_PORT"))
}
