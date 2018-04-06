package main

import (
	"github.com/catpie/musdk-go"
)

var (
	apiClient *musdk.Client
)

func InitWebApi() error {
	logger.Info("init mu api")
	apiClient = musdk.NewClient(Mu_Uri, Mu_Token, Mu_NodeID, musdk.TypeForward)
	apiClient.SetLogger(logger)
	go apiClient.UpdateTrafficDaemon()
	return nil
}
