package main

import (
	"log"
	"os"
	"smart-serve/models"
	"smart-serve/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// @title Smart Serve
// @description Smart Serve API
// @version 1.0
// @host localhost:5000
// @schemes http https
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {
	if os.Getenv("MODE") != "production" {
		err := godotenv.Load(".env")
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	models.ConnectDB()
	models.Migrate()

	r := gin.Default()

	routes.Config(r)

	r.Run(":" + os.Getenv("PORT"))
}
