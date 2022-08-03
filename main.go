package main

import (
	"aging-api/routes"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.GET("/api/v1", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"data": "Hello world"})
	})
	routes.AuthRoute(router)
	routes.MeasurementRoute(router)
	routes.SpiritRoute(router)
	routes.UserRoute(router)
	routes.VesselRoute(router)
	router.Run()
}
