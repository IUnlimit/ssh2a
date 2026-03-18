package rhttp

import (
	"fmt"
	"io/fs"
	"net"
	"net/http"
	"sync"

	"github.com/IUnlimit/ssh2a/internal/api"
	"github.com/IUnlimit/ssh2a/internal/cache"
	"github.com/IUnlimit/ssh2a/logger"
	"github.com/gin-gonic/gin"
	"github.com/libp2p/go-reuseport"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
)

func Listen(host string, port int, ipCache *cache.IPCache, webFS fs.FS, wg *sync.WaitGroup) {
	defer wg.Done()

	api.Cache = ipCache

	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = io.MultiWriter(os.Stdout, logger.Hook.GetWriter())

	engine := gin.Default()
	engine.Use(gin.Recovery())

	api.SetupRouter(engine, webFS)

	addr := fmt.Sprintf("%s:%d", host, port)
	log.Infof("HTTP server starting on %s", addr)

	listener, err := reuseport.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Failed to listen on HTTP port %s: %v", addr, err)
	}

	if err := http.Serve(listener, engine); err != nil {
		log.Fatalf("HTTP server error: %v", err)
	}
}

func multipleAbleHttpListen(host string, port int) net.Listener {
	httpListen, err := reuseport.Listen("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		log.Fatalf("Error occurred when listening resuse-http port(%d) on host(%s), %v", port, host, err)
	}
	return httpListen
}
