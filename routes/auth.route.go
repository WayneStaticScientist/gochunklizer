package routes

import (
	"github.com/doquangtan/socketio/v4"
	"github.com/gofiber/fiber/v2"
	"opechains.shop/chunklizer/v2/chunkedupload"
	"opechains.shop/chunklizer/v2/websocket"
)

func HandleRoutes(app *fiber.App) {
	io := socketio.New()
	wManager := websocket.NewWebSocketManager(io)
	chunkUpload := chunkedupload.InitChunkUploader(wManager)
	app.Put("/v1/chunk/update", chunkUpload.Update)
	app.Post("/v1/chunk/upload", chunkUpload.Upload)
	app.Delete("/v1/chunk/delete", chunkUpload.RequestDeleteFile)
	io.OnConnection(wManager.NewUserConnection)
	app.Use("/", io.FiberMiddleware) // handles WebSocket upgrade headers
	app.Route("/socket.io", io.FiberRoute)
	go chunkUpload.Work()
	go chunkUpload.CleanUpTemp()
}
