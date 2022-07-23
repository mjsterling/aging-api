package controllers

import (
	"aging-api/configs"
	"aging-api/models"
	"aging-api/responses"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var measurementCollection *mongo.Collection = configs.GetCollection(configs.DB, "measurements")
var validateMeasurement = validator.New()

func CreateMeasurement() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var measurement models.Measurement
		defer cancel()

		if err := c.BindJSON(&measurement); err != nil {
			c.JSON(http.StatusBadRequest,
				responses.Response{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}},
			)
			return
		}

		if validationErr := validateMeasurement.Struct(&measurement); validationErr != nil {
			c.JSON(http.StatusBadRequest, responses.Response{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}})
			return
		}

		newMeasurement := models.Measurement{
			Id:         primitive.NewObjectID(),
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
			c.JSON(http.StatusInternalServerError, responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		c.JSON(http.StatusCreated, responses.Response{Status: http.StatusCreated, Message: "success", Data: map[string]interface{}{"data": result.InsertedID}})
		return
	}
}

func GetMeasurement() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		measurementId := c.Param("id")
		var measurement models.Measurement
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(measurementId)

		err := measurementCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&measurement)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		c.JSON(http.StatusOK, responses.Response{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": measurement}})
		return
	}
}

func UpdateMeasurement() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		measurementId := c.Param("id")
		var measurement models.Measurement
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(measurementId)

		if err := c.BindJSON(&measurement); err != nil {
			c.JSON(http.StatusBadRequest, responses.Response{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		if validationErr := validateMeasurement.Struct(&measurement); validationErr != nil {
			c.JSON(http.StatusBadRequest, responses.Response{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}})
			return
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
			c.JSON(http.StatusInternalServerError, responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		var updatedMeasurement models.Measurement

		if result.MatchedCount == 1 {
			err := measurementCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&updatedMeasurement)

			if err != nil {
				c.JSON(http.StatusInternalServerError, responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
				return
			}
		}

		c.JSON(http.StatusOK, responses.Response{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": updatedMeasurement}})
		return
	}
}

func DeleteMeasurement() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		measurementId := c.Param("id")
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(measurementId)

		result, err := measurementCollection.DeleteOne(ctx, bson.M{"_id": objId})

		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		if result.DeletedCount < 1 {
			c.JSON(http.StatusNotFound, responses.Response{Status: http.StatusNotFound, Message: "error", Data: map[string]interface{}{"data": "measurement not found"}})
			return
		}

		c.JSON(http.StatusOK,
			responses.Response{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": "measurement deleted"}},
		)
		return
	}
}
