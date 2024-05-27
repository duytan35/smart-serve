package main

import (
	"smart-serve/models"
	"smart-serve/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	models.ConnectDB()

	r := gin.Default()

	routes.Config(r)

	r.Run(":5000")
}
