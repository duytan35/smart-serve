package main

import (
	"log"
	"os"
	"smart-serve/models"
	"smart-serve/routes"
	"smart-serve/utils"
	"smart-serve/validators"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
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

func configApp(r *gin.Engine) {
	models.ConnectDB()
	models.Migrate()
	utils.InitS3Uploader()
	utils.InitWebSocketServer(r)

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		if err := validators.RegisterCustomValidations(v); err != nil {
			panic(err)
		}
	}

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowHeaders = []string{"Origin", "Content-Type", "User-Agent", "Authorization", "Accept", "Cache-Control", "Pragma"}

	r.Use(cors.New(config))

	docs.SwaggerInfo.Host = os.Getenv("SERVER_URL")

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

	r.Run(":" + os.Getenv("PORT"))
}
