package api

import (
	"net/http"

	"github.com/IUnlimit/ssh2a/internal/db"
	"github.com/gin-gonic/gin"
)

func statusHandler(c *gin.Context) {
	ip := c.ClientIP()
	verified := Cache.IsVerified(ip)
	hasRecord := db.HasSSHAttempt(ip)
	hasLocalRecord := Cache.HasRecord(ip)

	c.JSON(http.StatusOK, gin.H{
		"ip":         ip,
		"verified":   verified,
		"has_record": hasRecord || hasLocalRecord,
	})
}
