package rhttp

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type LoginRequest struct {
	Password string `json:"password"`
	RemoteIP string `json:"remoteIP"`
}

func auth(c *gin.Context) {
	header := c.GetHeader("Authorization")
	if header != "" {
		authorization(header, c.RemoteIP(), c)
		return
	}

	var request LoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	authorization(request.Password, request.RemoteIP, c)
}

func authorization(password string, remoteIP string, c *gin.Context) {
	if password != "666" {
		log.Warnf("Worng password(%s) tried by ip(%s)", password, remoteIP)
		c.JSON(http.StatusForbidden, gin.H{
			"message": "Wrong password",
		})
		err := IPCache.UpdateIPStatus(remoteIP, false)
		if err != nil {
			log.Errorf("Error updating ip cache, %v", err)
			return
		}
		return
	}

	err := IPCache.UpdateIPStatus(remoteIP, true)
	if err != nil {
		log.Errorf("Error updating ip cache, %v", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{})
	return
}
