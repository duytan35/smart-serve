package routes

import (
	"net/http"
	_ "smart-serve/docs"

	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Config(r *gin.Engine) {
	v1 := r.Group("/api/v1")

	addAuthRoutes(v1)
	addRestaurantRoutes(v1)
	addAdminRoutes(v1)
	addFileRoutes(v1)
	addDishGroupRoutes(v1)
	addDishRoutes(v1)
	addTableRoutes(v1)
	addOrderRoutes(v1)

	v1.GET("/docs", func(c *gin.Context) { c.Redirect(http.StatusFound, "./docs/index.html") })
	v1.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
