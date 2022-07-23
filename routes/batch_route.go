package routes

import (
	"aging-api/controllers"

	"github.com/gin-gonic/gin"
)

func BatchRoute(router *gin.Engine) {
	router.GET("/api/v1/batches", controllers.GetBatch())
	router.POST("/api/v1/batches", controllers.CreateBatch())
	router.PUT("/api/v1/batches/:id", controllers.UpdateBatch())
	router.DELETE("/api/v1/batches/:id", controllers.DeleteBatch())
}
