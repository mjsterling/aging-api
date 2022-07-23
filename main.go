package main

import (
	"aging-api/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	routes.AuthRoute(router)
	routes.MeasurementRoute(router)
	routes.SpiritRoute(router)
	routes.UserRoute(router)
	routes.VesselRoute(router)
	router.Run()
}
