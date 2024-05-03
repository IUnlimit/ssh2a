package http

import (
	"fmt"
	"github.com/IUnlimit/ssh2a/internal/cache"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func Init(port int) {
	engine := gin.Default()
	engine.LoadHTMLGlob("./web/*")
	initRouter(engine)
	log.Infof("Starting http server on port: %d", port)
	err := engine.Run(fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("Http server exists with error: %v", err)
	}
}

func initRouter(engine *gin.Engine) {
	engine.NoRoute(func(ctx *gin.Context) {
		_ = ctx.AbortWithError(http.StatusBadGateway, http.ErrAbortHandler)
	})

	router := engine.Group("/")
	router.GET("/", middleIpAssert, handleAuth)

	api := engine.Group("/api")
	api.POST("/auth", middleIpAssert, handleAPIAuth)
}

func middleIpAssert(ctx *gin.Context) {
	ip := ctx.RemoteIP()
	if forwardIPs := ctx.Request.Header["X-Forwarded-For"]; forwardIPs != nil {
		log.Debugf("RemoteIP %s request with X-Forwarded-For %s", ip, forwardIPs)
		ip = forwardIPs[0]
	}

	if !cache.ContainsTempCache(ip) {
		_ = ctx.AbortWithError(http.StatusForbidden, http.ErrAbortHandler)
		return
	}
}
