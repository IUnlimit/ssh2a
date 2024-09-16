package rhttp

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func pass(c *gin.Context) {
	if !IPCache.CheckPassed(c.RemoteIP()) {
		c.String(http.StatusForbidden, "Access denied")
		return
	}
	c.HTML(http.StatusOK, "pass.tmpl", gin.H{
		"title":   "Login Success",
		"message": "Verification is successful, now you can directly connect via ssh",
	})
}
