package rhttp

import (
	"fmt"
	"github.com/IUnlimit/ssh2a/cache"
	"github.com/IUnlimit/ssh2a/logger"
	"github.com/gin-gonic/gin"
	"github.com/libp2p/go-reuseport"
	log "github.com/sirupsen/logrus"
	"io"
	"net"
	"net/http"
	"os"
	"sync"
	"time"
)

var IPCache = cache.NewIPCountCache("http", 256, 3*time.Hour)

func Listen(host string, port int, wg *sync.WaitGroup) {
	defer wg.Done()
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = io.MultiWriter(os.Stdout, logger.Hook.GetWriter())

	engine := gin.Default()
	engine.Use(gin.Recovery())
	engine.LoadHTMLGlob("templates/*")
	engine.GET("/", login)
	engine.GET("/pass", pass)

	v1 := engine.Group("/api/v1")
	v1.POST("/auth", auth)

	log.Infof("Http server starting on %s:%d", host, port)
	err := http.Serve(multipleAbleHttpListen(host, port), engine)
	if err != nil {
		log.Fatalf("Http server occurred error, %v", err)
	}
}

func multipleAbleHttpListen(host string, port int) net.Listener {
	httpListen, err := reuseport.Listen("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		log.Fatalf("Error occurred when listening resuse-http port(%d) on host(%s), %v", port, host, err)
	}
	return httpListen
}
