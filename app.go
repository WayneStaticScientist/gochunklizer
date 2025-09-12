package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	"opechains.shop/chunklizer/v2/routes"
)

func main() {
	app := fiber.New()
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	app.Use(cors.New(cors.Config{
		AllowHeaders: "Origin, Content-Type, Accept, X-Object-Type, X-Upload-Id,X-Object-Id,X-Object-Key",
	}))

	app.Use(logger.New())
	routes.HandleRoutes(app)
	app.Listen(":5142")
}
