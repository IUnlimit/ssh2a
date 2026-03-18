package api

import (
	"io/fs"
	"net/http"

	"github.com/gin-gonic/gin"
)

// SetupRouter 注册所有路由
func SetupRouter(engine *gin.Engine, webFS fs.FS) {
	// API 路由
	v1 := engine.Group("/api/v1")
	{
		v1.POST("/auth", authHandler)
		v1.GET("/status", statusHandler)

		admin := v1.Group("/admin")
		{
			admin.GET("/honeypot", adminHoneypotHandler)
			admin.GET("/rejected", adminRejectedHandler)
			admin.GET("/verified", adminVerifiedHandler)
			admin.GET("/stats", adminStatsHandler)
		}
	}

	// 静态文件服务（Vue SPA）
	if webFS != nil {
		engine.NoRoute(func(c *gin.Context) {
			// 尝试提供静态文件
			path := c.Request.URL.Path
			if path == "/" {
				path = "/index.html"
			}
			c.FileFromFS(path, http.FS(webFS))
		})
	}
}
