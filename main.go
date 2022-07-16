package main

import (
	"aging-app/configs"
	"aging-app/routes"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()
	configs.ConnectDB()

	routes.UserRoute(app)

	app.Listen(":6000")
}
