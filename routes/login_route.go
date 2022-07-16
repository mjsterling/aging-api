package routes

import (
	"aging-app/controllers"

	"github.com/gofiber/fiber/v2"
)

func LoginRoute(app *fiber.App) {
	app.Post("/login", controllers.Login)
}
