package main

import (
	"github.com/IUnlimit/ssh2a/cmd/ssh2a"
	"github.com/IUnlimit/ssh2a/internal/conf"
	"github.com/IUnlimit/ssh2a/internal/logger"
)

func main() {
	conf.Init()
	logger.Init()
	ssh2a.Run()
}
