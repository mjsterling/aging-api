package routes

import (
	"aging-api/controllers"

	"github.com/gin-gonic/gin"
)

func SpiritRoute(router *gin.Engine) {
	router.GET("/api/v1/spirits", controllers.GetSpirit())
	router.POST("/api/v1/spirits", controllers.CreateSpirit())
	router.PUT("/api/v1/spirits/:id", controllers.UpdateSpirit())
	router.DELETE("/api/v1/spirits/:id", controllers.DeleteSpirit())
}
