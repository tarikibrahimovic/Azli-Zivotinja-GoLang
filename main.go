package main

import (
	"azil-app/configs"
	"azil-app/routes"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	//run database
	configs.ConnectDB()

	//routes
	routes.AnimalUserRoute(app) //add this

	app.Listen(":6000")
}
