
package routes


import (
	"c/handlers"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine){
	api := router.Group("/api")
	{
		api.POST("/shorten", handlers.CreateShortURL)
		api.GET("/stats/:code", handlers.GetStats)
		api.DELETE("/:code", handlers.DeleteShortURL)
	}

	// Redirect route (root level): /:code
	router.GET(":code", handlers.RedirectShortURL)
}