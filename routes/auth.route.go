package routes

import "github.com/gofiber/fiber/v2"

func HandleRoutes(app *fiber.App) {
	app.Post("/v1/chuck/initiate")
	app.Post("/v1/chuck/upload")
}
