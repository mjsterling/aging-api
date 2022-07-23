package routes

import (
	"aging-api/controllers"

	"github.com/gin-gonic/gin"
)

func UserRoute(router *gin.Engine) {
	router.GET("/api/v1/users", controllers.GetUser())
	router.POST("/api/v1/users", controllers.CreateUser())
	router.PUT("/api/v1/users/:id", controllers.UpdateUser())
	router.DELETE("/api/v1/users/:id", controllers.DeleteUser())
}
