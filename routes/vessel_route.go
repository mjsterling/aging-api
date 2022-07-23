package routes

import (
	"aging-api/controllers"

	"github.com/gin-gonic/gin"
)

func VesselRoute(router *gin.Engine) {
	router.GET("/api/v1/vessels", controllers.GetVessel())
	router.POST("/api/v1/vessels", controllers.CreateVessel())
	router.PUT("/api/v1/vessels/:id", controllers.UpdateVessel())
	router.DELETE("/api/v1/vessels/:id", controllers.DeleteVessel())
}
