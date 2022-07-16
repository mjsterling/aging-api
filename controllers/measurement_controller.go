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

var measurementCollection *mongo.Collection = configs.GetCollection(configs.DB, "measurements")
var validateMeasurement = validator.New()

func CreateMeasurement(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var measurement models.Measurement
	defer cancel()

	if err := c.BodyParser(&measurement); err != nil {
		return c.Status(http.StatusBadRequest).JSON(
			responses.Response{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}},
		)
	}

	if validationErr := validateMeasurement.Struct(&measurement); validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.Response{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": validationErr.Error()}})
	}

	newMeasurement := models.Measurement{
		ID:         primitive.NewObjectID(),
		CreatedAt:  primitive.NewDateTimeFromTime(time.Now()),
		ABV:        measurement.ABV,
		Image:      measurement.Image,
		Nose:       measurement.Nose,
		ForePalate: measurement.ForePalate,
		MidPalate:  measurement.MidPalate,
		Finish:     measurement.Finish,
		Notes:      measurement.Notes,
	}

	result, err := measurementCollection.InsertOne(ctx, newMeasurement)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	return c.Status(http.StatusCreated).JSON(responses.Response{Status: http.StatusCreated, Message: "success", Data: &fiber.Map{"data": result.InsertedID}})
}

func GetMeasurement(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	measurementId := c.Params("id")
	var measurement models.Measurement
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(measurementId)

	err := measurementCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&measurement)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	return c.Status(http.StatusOK).JSON(responses.Response{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": measurement}})
}

func UpdateMeasurement(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	measurementId := c.Params("id")
	var measurement models.Measurement
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(measurementId)

	if err := c.BodyParser(&measurement); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.Response{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	if validationErr := validateMeasurement.Struct(&measurement); validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.Response{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": validationErr.Error()}})
	}

	update := bson.M{
		"ABV":        measurement.ABV,
		"Image":      measurement.Image,
		"Nose":       measurement.Nose,
		"ForePalate": measurement.ForePalate,
		"MidPalate":  measurement.MidPalate,
		"Finish":     measurement.Finish,
		"Notes":      measurement.Notes,
	}

	result, err := measurementCollection.UpdateOne(ctx, bson.M{"_id": objId}, bson.M{"$set": update})

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	var updatedMeasurement models.Measurement

	if result.MatchedCount == 1 {
		err := measurementCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&updatedMeasurement)

		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
		}
	}

	return c.Status(http.StatusOK).JSON(responses.Response{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": updatedMeasurement}})

}

func DeleteMeasurement(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	measurementId := c.Params("id")
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(measurementId)

	result, err := measurementCollection.DeleteOne(ctx, bson.M{"_id": objId})

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	if result.DeletedCount < 1 {
		return c.Status(http.StatusNotFound).JSON(responses.Response{Status: http.StatusNotFound, Message: "error", Data: &fiber.Map{"data": "measurement not found"}})
	}

	return c.Status(http.StatusOK).JSON(
		responses.Response{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": "measurement deleted"}},
	)
}
