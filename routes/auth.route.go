package routes

import (
	"github.com/gofiber/fiber/v2"
	"opechains.shop/chunklizer/v2/chunkedupload"
)

func HandleRoutes(app *fiber.App) {
	chunkUpload := chunkedupload.InitChunkUploader()
	app.Put("/v1/chunk/update", chunkUpload.Update)
	app.Post("/v1/chunk/upload", chunkUpload.Upload)
	app.Delete("/v1/chunk/delete", chunkUpload.RequestDeleteFile)
	go chunkUpload.Work()
	go chunkUpload.CleanUpTemp()
}
