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

func CreateVessel() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var vessel models.Vessel
		defer cancel()

		if err := c.BindJSON(&vessel); err != nil {
			c.JSON(http.StatusBadRequest,
				responses.Response{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}},
			)
		}

		if validationErr := validateVessel.Struct(&vessel); validationErr != nil {
			c.JSON(http.StatusBadRequest, responses.Response{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}})
			return
		}

		newVessel := models.Vessel{
			Id:        primitive.NewObjectID(),
			Batches:   vessel.Batches,
			CreatedAt: primitive.NewDateTimeFromTime(time.Now()),
			Volume:    vessel.Volume,
			Material:  vessel.Material,
			Process:   vessel.Process,
		}

		result, err := vesselCollection.InsertOne(ctx, newVessel)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		c.JSON(http.StatusCreated, responses.Response{Status: http.StatusCreated, Message: "success", Data: map[string]interface{}{"data": result.InsertedID}})
		return
	}
}

func GetVessel() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		vesselId := c.Param("id")
		var vessel models.Vessel
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(vesselId)

		err := vesselCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&vessel)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		c.JSON(http.StatusOK, responses.Response{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": vessel}})
		return
	}
}

func UpdateVessel() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		vesselId := c.Param("id")
		var vessel models.Vessel
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(vesselId)

		if err := c.BindJSON(&vessel); err != nil {
			c.JSON(http.StatusBadRequest, responses.Response{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		if validationErr := validateVessel.Struct(&vessel); validationErr != nil {
			c.JSON(http.StatusBadRequest, responses.Response{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}})
			return
		}

		update := bson.M{
			"Volume":   vessel.Volume,
			"Material": vessel.Material,
			"Process":  vessel.Process,
		}

		result, err := vesselCollection.UpdateOne(ctx, bson.M{"_id": objId}, bson.M{"$set": update})

		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		var updatedVessel models.Vessel

		if result.MatchedCount == 1 {
			err := vesselCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&updatedVessel)

			if err != nil {
				c.JSON(http.StatusInternalServerError, responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
				return
			}
		}

		c.JSON(http.StatusOK, responses.Response{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": updatedVessel}})
		return

	}
}

func DeleteVessel() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		vesselId := c.Param("id")
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(vesselId)

		result, err := vesselCollection.DeleteOne(ctx, bson.M{"_id": objId})

		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		if result.DeletedCount < 1 {
			c.JSON(http.StatusNotFound, responses.Response{Status: http.StatusNotFound, Message: "error", Data: map[string]interface{}{"data": "vessel not found"}})
			return
		}

		c.JSON(http.StatusOK,
			responses.Response{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": "vessel deleted"}},
		)
		return
	}
}
