package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"opechains.shop/chunklizer/v2/routes"
)

func main() {
	app := fiber.New()

	app.Use(cors.New())
	app.Use(logger.New())
	routes.HandleRoutes(app)
	app.Listen(":5142")
}
