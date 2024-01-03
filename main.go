package main

import (
	"azil-app/configs"
	"azil-app/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	app := fiber.New()

	app.Use(cors.New())

	configs.ConnectDB()

	routes.AnimalUserRoute(app)

	app.Listen(":6000")
}
