package routes

import (
	"github.com/bielgennaro/vehicle-vision-api/api/controllers"
	"github.com/bielgennaro/vehicle-vision-api/api/middlewares"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	router.Use(middlewares.CORSMiddleware())

	router.GET("/health", controllers.HealthCheck)
}
