package main

import (
	"github.com/gofiber/fiber/v2"
	"opechains.shop/chunklizer/v2/routes"
)

func main() {
	app := fiber.New()
	routes.HandleRoutes(app)
	app.Listen(":3000")
}
