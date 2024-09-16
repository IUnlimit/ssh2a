package main

import (
	"github.com/IUnlimit/ssh2a/cmd/rhttp"
	"github.com/IUnlimit/ssh2a/cmd/rssh"
	"github.com/IUnlimit/ssh2a/conf"
	"github.com/IUnlimit/ssh2a/logger"
	"github.com/IUnlimit/ssh2a/tools"
	log "github.com/sirupsen/logrus"
	"sync"
)

func main() {
	logger.Init()
	initAuth()
	bind := conf.Config.Bind
	var wg sync.WaitGroup
	wg.Add(1)
	go rhttp.Listen(bind.Host, bind.HttpPort, &wg)
	wg.Add(1)
	go rssh.Listen(bind.Host, bind.SSHPort, &wg)
	wg.Wait()
	log.Info("All Goroutines have finished, main thread exiting")
}

func initAuth() {
	auth := conf.Config.Authorization
	log.Infof("Authorization type: %s", auth.Type)
	if auth.Type == "basic" {
	} else if auth.Type == "authenticator" {
		err := tools.PrintQRCode(auth.Authenticator.PrivateSecret)
		if err != nil {
			log.Errorf("Failed to print QRCode, if your secret not changed, ignore it, %v", err)
		}
	} else {
		log.Fatalf("Unknown authorization type: %s", auth.Type)
	}
}
