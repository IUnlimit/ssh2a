package main

import (
	"github.com/IUnlimit/ssh2a/cmd/rhttp"
	"github.com/IUnlimit/ssh2a/cmd/rssh"
	"github.com/IUnlimit/ssh2a/logger"
	log "github.com/sirupsen/logrus"
	"sync"
)

func main() {
	logger.Init()
	log.Info("Loading config")
	host := "0.0.0.0"
	port := 9022
	var wg sync.WaitGroup
	wg.Add(1)
	go rhttp.Listen(host, port, &wg)
	wg.Add(1)
	go rssh.Listen(host, port, &wg)
	wg.Wait()
	log.Info("All Goroutines have finished, main thread exiting")
}
