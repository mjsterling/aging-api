package controllers

import (
	"aging-app/configs"
	"aging-app/models"
	"aging-app/responses"
	"context"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var spiritCollection *mongo.Collection = configs.GetCollection(configs.DB, "spirits")
var validateSpirit = validator.New()

func CreateSpirit(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var spirit models.Spirit
	defer cancel()

	if err := c.BodyParser(&spirit); err != nil {
		return c.Status(http.StatusBadRequest).JSON(
			responses.Response{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}},
		)
	}

	if validationErr := validateSpirit.Struct(&spirit); validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.Response{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": validationErr.Error()}})
	}

	newSpirit := models.Spirit{
		ID:         primitive.NewObjectID(),
		Batches:    spirit.Batches,
		CreatedAt:  primitive.NewDateTimeFromTime(time.Now()),
		Volume:     spirit.Volume,
		Name:       spirit.Name,
		Type:       spirit.Type,
		InitialABV: spirit.InitialABV,
		RecipeName: spirit.RecipeName,
	}

	result, err := spiritCollection.InsertOne(ctx, newSpirit)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	return c.Status(http.StatusCreated).JSON(responses.Response{Status: http.StatusCreated, Message: "success", Data: &fiber.Map{"data": result.InsertedID}})
}

func GetSpirit(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	spiritId := c.Params("id")
	var spirit models.Spirit
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(spiritId)

	err := spiritCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&spirit)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	return c.Status(http.StatusOK).JSON(responses.Response{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": spirit}})
}

func UpdateSpirit(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	spiritId := c.Params("id")
	var spirit models.Spirit
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(spiritId)

	if err := c.BodyParser(&spirit); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.Response{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	if validationErr := validateSpirit.Struct(&spirit); validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.Response{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": validationErr.Error()}})
	}

	update := bson.M{
		"Volume":     spirit.Volume,
		"Name":       spirit.Name,
		"Type":       spirit.Type,
		"InitialABV": spirit.InitialABV,
		"RecipeName": spirit.RecipeName,
	}

	result, err := spiritCollection.UpdateOne(ctx, bson.M{"_id": objId}, bson.M{"$set": update})

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	var updatedSpirit models.Spirit

	if result.MatchedCount == 1 {
		err := spiritCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&updatedSpirit)

		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
		}
	}

	return c.Status(http.StatusOK).JSON(responses.Response{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": updatedSpirit}})

}

func DeleteSpirit(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	spiritId := c.Params("id")
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(spiritId)

	result, err := spiritCollection.DeleteOne(ctx, bson.M{"_id": objId})

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	if result.DeletedCount < 1 {
		return c.Status(http.StatusNotFound).JSON(responses.Response{Status: http.StatusNotFound, Message: "error", Data: &fiber.Map{"data": "spirit not found"}})
	}

	return c.Status(http.StatusOK).JSON(
		responses.Response{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": "spirit deleted"}},
	)
}
