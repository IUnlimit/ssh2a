package api

import (
	"io/fs"
	"net/http"
	"strings"

	"github.com/IUnlimit/ssh2a/conf"
	"github.com/gin-gonic/gin"
)

// SetupRouter 注册所有路由
func SetupRouter(engine *gin.Engine, webFS fs.FS) {
	// API 路由
	v1 := engine.Group("/api/v1")
	{
		v1.POST("/auth", authHandler)
		v1.GET("/status", statusHandler)

		admin := v1.Group("/admin", adminIPGuard())
		{
			admin.GET("/honeypot", adminHoneypotHandler)
			admin.GET("/rejected", adminRejectedHandler)
			admin.GET("/verified", adminVerifiedHandler)
			admin.GET("/stats", adminStatsHandler)
		}
	}

	// 静态文件服务（Vue SPA）
	if webFS != nil {
		httpFS := http.FS(webFS)
		engine.NoRoute(func(c *gin.Context) {
			path := c.Request.URL.Path

			// 不处理 API 路径的 404
			if strings.HasPrefix(path, "/api/") {
				c.JSON(http.StatusNotFound, gin.H{"message": "not found"})
				return
			}

			// 尝试打开请求的文件
			f, err := webFS.Open(strings.TrimPrefix(path, "/"))
			if err == nil {
				f.Close()
				// 文件存在，直接提供
				http.FileServer(httpFS).ServeHTTP(c.Writer, c.Request)
				return
			}

			// 文件不存在，SPA fallback 到 index.html
			c.Request.URL.Path = "/"
			http.FileServer(httpFS).ServeHTTP(c.Writer, c.Request)
		})
	}
}

// adminIPGuard 管理台 IP 白名单中间件
func adminIPGuard() gin.HandlerFunc {
	return func(c *gin.Context) {
		cfg := conf.Config.Admin
		if cfg == nil || len(cfg.AllowedHosts) == 0 {
			c.Next()
			return
		}

		clientIP := c.ClientIP()
		for _, host := range cfg.AllowedHosts {
			if clientIP == host {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"message": "access denied",
		})
	}
}
