package main

import (
	"log"
	"os"
	"smart-serve/models"
	"smart-serve/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	models.ConnectDB()
	models.Migrate()

	r := gin.Default()

	routes.Config(r)

	r.Run(":" + os.Getenv("PORT"))
}
