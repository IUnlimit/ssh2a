package api

import (
	"net/http"
	"strconv"

	"github.com/IUnlimit/ssh2a/internal/db"
	"github.com/gin-gonic/gin"
)

func parsePagination(c *gin.Context) (int, int) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	return page, pageSize
}

func adminHoneypotHandler(c *gin.Context) {
	page, pageSize := parsePagination(c)
	data, total, err := db.GetHoneypotCredentials(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data":  data,
		"total": total,
		"page":  page,
	})
}

func adminRejectedHandler(c *gin.Context) {
	page, pageSize := parsePagination(c)
	data, total, err := db.GetRejectedIPs(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data":  data,
		"total": total,
		"page":  page,
	})
}

func adminVerifiedHandler(c *gin.Context) {
	page, pageSize := parsePagination(c)
	data, total, err := db.GetVerifiedIPs(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data":  data,
		"total": total,
		"page":  page,
	})
}

func adminStatsHandler(c *gin.Context) {
	stats, err := db.GetStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, stats)
}
