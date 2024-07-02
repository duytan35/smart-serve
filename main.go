package main

import (
	"log"
	"os"
	"smart-serve/models"
	"smart-serve/routes"
	"smart-serve/utils"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"smart-serve/docs"
)

// @title Smart Serve
// @description Smart Serve API
// @version 1.0
// @schemes http https
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func configSwaggerHost() string {
	if os.Getenv("MODE") == "release" {
		return "34.126.68.84:5000"
	}
	return "localhost:5000"
}

func configApp(r *gin.Engine) {
	models.ConnectDB()
	models.Migrate()
	utils.InitS3Uploader()

	config := cors.DefaultConfig()
	config.AllowCredentials = true
	config.AllowWildcard = true
	config.AllowOrigins = []string{"http://*", "https://*"}

	r.Use(cors.New(config))

	docs.SwaggerInfo.Host = configSwaggerHost()

	routes.Config(r)
}

func main() {
	if os.Getenv("MODE") != "release" {
		err := godotenv.Load(".env")
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	gin.SetMode(os.Getenv("MODE"))

	r := gin.Default()

	configApp(r)
	configApp(r)

	r.Run(":" + os.Getenv("PORT"))
}
