package main

import (
	"log"

	"github.com/bielgennaro/vehicle-vision-api/api/routes"
	"github.com/bielgennaro/vehicle-vision-api/configs"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Println(".env file not found or not loaded")
	}

	config := configs.LoadConfig()

	if config == nil {
		log.Fatal("Failed to load configuration")
	}

	if config.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	routes.SetupRoutes(router)

	log.Printf("🧟 Server running on port:  %s\n", config.Port)
	err = router.Run(":" + config.Port)
	if err != nil {
		log.Fatalf("😢 Failed to start server: %v", err)
	}
}
