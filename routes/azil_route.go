package routes

import (
	"azil-app/controllers"
	"github.com/gofiber/fiber/v2"
)

func AnimalUserRoute(app *fiber.App) {
	app.Post("/animal_user", controllers.CreateAnimalUser)
	app.Get("/animal_users", controllers.GetAllAnimalUsers)
}
