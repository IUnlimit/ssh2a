package rhttp

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func login(c *gin.Context) {
	if IPCache.CheckPassed(c.RemoteIP()) {
		c.Redirect(http.StatusOK, "/pass")
		return
	}
	c.HTML(http.StatusOK, "login.tmpl", gin.H{
		"remoteIP": c.RemoteIP(),
	})
}
