package routes

import (
	"github.com/gofiber/fiber/v2"
	"opechains.shop/chunklizer/v2/chunkedupload"
)

func HandleRoutes(app *fiber.App) {
	chunkUpload := chunkedupload.InitChunkUploader()
	app.Post("/v1/chuck/upload", chunkUpload.Upload)
}
