package main

import (
	"io/fs"
	"sync"

	"github.com/IUnlimit/ssh2a/cmd/rhttp"
	"github.com/IUnlimit/ssh2a/cmd/rssh"
	"github.com/IUnlimit/ssh2a/conf"
	"github.com/IUnlimit/ssh2a/internal/cache"
	"github.com/IUnlimit/ssh2a/internal/db"
	"github.com/IUnlimit/ssh2a/internal/honeypot"
	"github.com/IUnlimit/ssh2a/logger"
	"github.com/IUnlimit/ssh2a/tools"
	log "github.com/sirupsen/logrus"
)

func main() {
	logger.Init()
	initAuth()

	// 初始化数据库
	db.Init(conf.Config.Database)

	// 初始化蜜罐（加载本机 sshd host key）
	honeypot.Init()

	// 初始化 IP 缓存
	ipCache := cache.NewIPCache(
		conf.Config.Honeypot.TriggerTimeout,
		conf.Config.Auth.ValidDuration,
	)

	// 获取嵌入的前端资源
	webFS, err := fs.Sub(webDist, "web/dist")
	if err != nil {
		log.Warnf("Failed to load embedded web assets: %v", err)
		webFS = nil
	}

	bind := conf.Config.Bind
	var wg sync.WaitGroup

	wg.Add(1)
	go rhttp.Listen(bind.Host, bind.HttpPort, ipCache, webFS, &wg)

	wg.Add(1)
	go rssh.Listen(bind.Host, bind.SSHPort, ipCache, &wg)

	wg.Wait()
	log.Info("All goroutines have finished, main thread exiting")
}

func initAuth() {
	auth := conf.Config.Authorization
	log.Infof("Authorization type: %s", auth.Type)
	if auth.Type == "basic" {
		// no-op
	} else if auth.Type == "authenticator" {
		err := tools.PrintQRCode(auth.Authenticator.PrivateSecret)
		if err != nil {
			log.Errorf("Failed to print QRCode, if your secret not changed, ignore it, %v", err)
		}
	} else {
		log.Fatalf("Unknown authorization type: %s", auth.Type)
	}
}
