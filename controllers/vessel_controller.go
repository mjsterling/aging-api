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

var vesselCollection *mongo.Collection = configs.GetCollection(configs.DB, "vessels")
var validateVessel = validator.New()
var _ = validateVessel.RegisterValidation(
	"material",
	func(fl validator.FieldLevel) bool {
		if fl.Field().String() == "French Oak" {
			return true
		}
		if fl.Field().String() == "American Oak" {
			return true
		}
		if fl.Field().String() == "Stainless" {
			return true
		}
		if fl.Field().String() == "Glass" {
			return true
		}

		return false
	},
)
var _ = validateVessel.RegisterValidation(
	"process",
	func(fl validator.FieldLevel) bool {
		if fl.Field().String() == "Charred" {
			return true
		}
		if fl.Field().String() == "Toasted" {
			return true
		}
		if fl.Field().String() == "" {
			return true
		}

		return false
	},
)

func CreateVessel(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var vessel models.Vessel
	defer cancel()

	if err := c.BodyParser(&vessel); err != nil {
		return c.Status(http.StatusBadRequest).JSON(
			responses.Response{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}},
		)
	}

	if validationErr := validateVessel.Struct(&vessel); validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.Response{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": validationErr.Error()}})
	}

	newVessel := models.Vessel{
		ID:        primitive.NewObjectID(),
		Batches:   vessel.Batches,
		CreatedAt: primitive.NewDateTimeFromTime(time.Now()),
		Volume:    vessel.Volume,
		Material:  vessel.Material,
		Process:   vessel.Process,
	}

	result, err := vesselCollection.InsertOne(ctx, newVessel)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	return c.Status(http.StatusCreated).JSON(responses.Response{Status: http.StatusCreated, Message: "success", Data: &fiber.Map{"data": result.InsertedID}})
}

func GetVessel(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	vesselId := c.Params("id")
	var vessel models.Vessel
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(vesselId)

	err := vesselCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&vessel)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	return c.Status(http.StatusOK).JSON(responses.Response{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": vessel}})
}

func UpdateVessel(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	vesselId := c.Params("id")
	var vessel models.Vessel
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(vesselId)

	if err := c.BodyParser(&vessel); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.Response{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	if validationErr := validateVessel.Struct(&vessel); validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.Response{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": validationErr.Error()}})
	}

	update := bson.M{
		"Volume":   vessel.Volume,
		"Material": vessel.Material,
		"Process":  vessel.Process,
	}

	result, err := vesselCollection.UpdateOne(ctx, bson.M{"_id": objId}, bson.M{"$set": update})

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	var updatedVessel models.Vessel

	if result.MatchedCount == 1 {
		err := vesselCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&updatedVessel)

		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
		}
	}

	return c.Status(http.StatusOK).JSON(responses.Response{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": updatedVessel}})

}

func DeleteVessel(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	vesselId := c.Params("id")
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(vesselId)

	result, err := vesselCollection.DeleteOne(ctx, bson.M{"_id": objId})

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	if result.DeletedCount < 1 {
		return c.Status(http.StatusNotFound).JSON(responses.Response{Status: http.StatusNotFound, Message: "error", Data: &fiber.Map{"data": "vessel not found"}})
	}

	return c.Status(http.StatusOK).JSON(
		responses.Response{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": "vessel deleted"}},
	)
}
