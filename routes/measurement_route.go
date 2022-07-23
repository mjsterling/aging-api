package routes

import (
	"aging-api/controllers"

	"github.com/gin-gonic/gin"
)

func MeasurementRoute(router *gin.Engine) {
	router.GET("/api/v1/measurements", controllers.GetMeasurement())
	router.POST("/api/v1/measurements", controllers.CreateMeasurement())
	router.PUT("/api/v1/measurements/:id", controllers.UpdateMeasurement())
	router.DELETE("/api/v1/measurements/:id", controllers.DeleteMeasurement())
}
