package controllers

import (
	"azil-app/configs"
	"azil-app/models"
	"azil-app/responses"
	"context"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var animalUsersCollection *mongo.Collection = configs.GetCollection(configs.DB, "Azil_Animal-Users")
var validate = validator.New()

func CreateAnimalUser(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var animalUser models.Collection
	defer cancel()

	//validate the request body
	if err := c.BodyParser(&animalUser); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.AnimalUserResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	//use the validator library to validate required fields
	if validationErr := validate.Struct(&animalUser); validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.AnimalUserResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": validationErr.Error()}})
	}

	// Check if the user already has 10 or more animals
	count, err := animalUsersCollection.CountDocuments(ctx, bson.M{"user_id": animalUser.UserId})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.AnimalUserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}
	if count >= 10 {
		return c.Status(http.StatusBadRequest).JSON(responses.AnimalUserResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": "User cannot get more than 10 animals"}})
	}

	newAnimalUser := models.Collection{
		Id:           primitive.NewObjectID(),
		UserId:       animalUser.UserId,
		AnimalId:     animalUser.AnimalId,
		DateOfTaking: time.Now().UTC(),
	}

	result, err := animalUsersCollection.InsertOne(ctx, newAnimalUser)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.AnimalUserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	// Query the database for the newly inserted document
	var createdAnimalUser models.Collection
	err = animalUsersCollection.FindOne(ctx, bson.M{"_id": result.InsertedID}).Decode(&createdAnimalUser)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.AnimalUserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	return c.Status(http.StatusCreated).JSON(responses.AnimalUserResponse{Status: http.StatusCreated, Message: "success", Data: &fiber.Map{"data": createdAnimalUser}})
}

func GetAllAnimalUsers(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var animalUsers []models.Collection
	defer cancel()

	results, err := animalUsersCollection.Find(ctx, bson.M{})

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.AnimalUserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	//reading from the db in an optimal way
	defer results.Close(ctx)
	for results.Next(ctx) {
		var singleAnimalUser models.Collection
		if err = results.Decode(&singleAnimalUser); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.AnimalUserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
		}

		animalUsers = append(animalUsers, singleAnimalUser)
	}

	return c.Status(http.StatusOK).JSON(
		responses.AnimalUserResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": animalUsers}},
	)
}
