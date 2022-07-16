package controllers

import (
	"aging-app/auth"
	"aging-app/models"
	"aging-app/responses"
	"context"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

func Login(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var requestBody models.User
	var fetchedUser models.User
	defer cancel()

	if err := c.BodyParser(&requestBody); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.Response{Status: http.StatusBadRequest, Message: "Bad Request", Data: &fiber.Map{"data": err.Error()}})
	}

	if err := userCollection.FindOne(ctx, bson.M{"email": requestBody.Email}).Decode(&fetchedUser); err != nil {
		return c.Status(http.StatusNotFound).JSON(responses.Response{Status: http.StatusInternalServerError, Message: "Account not found", Data: &fiber.Map{"data": err.Error()}})
	}

	if authorized := auth.CheckPasswordHash(requestBody.Password, fetchedUser.Password); authorized == false {
		return c.Status(http.StatusUnauthorized).JSON(responses.Response{Status: http.StatusUnauthorized, Message: "error", Data: &fiber.Map{"data": "Incorrect password"}})
	}

	// create JWT token
	token, err := auth.CreateJWT(fetchedUser.ID.Hex())

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	return c.Status(http.StatusOK).JSON(
		responses.Response{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"token": token}},
	)

}
