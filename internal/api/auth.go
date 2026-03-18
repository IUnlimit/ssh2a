package api

import (
	"net/http"

	"github.com/IUnlimit/ssh2a/conf"
	"github.com/IUnlimit/ssh2a/internal/cache"
	"github.com/IUnlimit/ssh2a/internal/db"
	"github.com/IUnlimit/ssh2a/tools"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// Cache 全局 IP 缓存引用，由 main 注入
var Cache *cache.IPCache

type authRequest struct {
	Password string `json:"password"`
	RemoteIP string `json:"remoteIP"`
}

func authHandler(c *gin.Context) {
	// 优先检查 Authorization 请求头
	header := c.GetHeader("Authorization")
	if header != "" {
		doAuth(header, c.ClientIP(), c)
		return
	}

	var req authRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	ip := req.RemoteIP
	if ip == "" {
		ip = c.ClientIP()
	}
	doAuth(req.Password, ip, c)
}

func doAuth(password, ip string, c *gin.Context) {
	config := conf.Config.Authorization
	method := config.Type
	ok := checkAccess(password)

	// 记录到数据库
	db.RecordAuth(ip, method, ok)

	if !ok {
		log.Warnf("Wrong password or 2fa key(%s) tried by ip(%s)", password, ip)
		Cache.SetRejected(ip)
		c.JSON(http.StatusForbidden, gin.H{"message": "Wrong password or 2fa key"})
		return
	}

	Cache.SetVerified(ip)
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}

func checkAccess(password string) bool {
	config := conf.Config.Authorization
	if config.Type == "basic" {
		return password == config.Basic.Secret
	} else if config.Type == "authenticator" {
		totP, err := tools.TotP(config.Authenticator.PrivateSecret)
		if err != nil {
			log.Errorf("Failed to get 2fa key: %v", err)
			return false
		}
		return password == totP
	}
	return false
}
