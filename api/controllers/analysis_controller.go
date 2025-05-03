package controllers

import (
	"net/http"

	"github.com/bielgennaro/vehicle-vision-api/configs"
	"github.com/bielgennaro/vehicle-vision-api/models"
	"github.com/gin-gonic/gin"
)

// GetAllAnalyses retorna todas as análises
func GetAllAnalyses(c *gin.Context) {
	var analyses []models.Analysis

	result := configs.DB.Find(&analyses)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching analyses"})
		return
	}

	c.JSON(http.StatusOK, analyses)
}

// GetAnalysis retorna uma análise pelo ID
func GetAnalysis(c *gin.Context) {
	id := c.Param("id")

	var analysis models.Analysis
	result := configs.DB.First(&analysis, id)

	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Analysis not found"})
		return
	}

	c.JSON(http.StatusOK, analysis)
}
