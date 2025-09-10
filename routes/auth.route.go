package routes

import (
	"github.com/gofiber/fiber/v2"
	"opechains.shop/chunklizer/v2/chunkedupload"
	"opechains.shop/chunklizer/v2/database"
)

func HandleRoutes(app *fiber.App) {
	database := database.InitDatabase()
	chuckUploader := chunkedupload.InitChunkUploader(database)
	app.Post("/v1/chuck/initiate")
	app.Post("/v1/chuck/upload")
}
