package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":      "Fine",
		"uptime":      time.Since(time.Now().Add(-24 * time.Hour)).String(),
		"version":     "1.0.0",
		"environment": "development",
		"timestamp":   time.Now().Format(time.RFC3339),
		"service":     "vehicle-vision-api",
	})
}
