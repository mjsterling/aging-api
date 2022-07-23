package routes

import (
	"aging-api/controllers"

	"github.com/gin-gonic/gin"
)

func AuthRoute(router *gin.Engine) {
	router.POST("/api/v1/login", controllers.Login())
}
