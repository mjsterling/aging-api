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

var batchCollection *mongo.Collection = configs.GetCollection(configs.DB, "batches")
var validateBatch = validator.New()

func CreateBatch() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var batch models.Batch
		defer cancel()

		if err := c.BindJSON(&batch); err != nil {
			c.JSON(http.StatusBadRequest,
				responses.Response{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}},
			)
			return
		}

		if validationErr := validateBatch.Struct(&batch); validationErr != nil {
			c.JSON(http.StatusBadRequest, responses.Response{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}})
			return
		}

		newBatch := models.Batch{
			Id:           primitive.NewObjectID(),
			Vessels:      batch.Vessels,
			Measurements: batch.Measurements,
			CreatedAt:    primitive.NewDateTimeFromTime(time.Now()),
			Volume:       batch.Volume,
		}

		result, err := batchCollection.InsertOne(ctx, newBatch)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		c.JSON(http.StatusCreated, responses.Response{Status: http.StatusCreated, Message: "success", Data: map[string]interface{}{"data": result.InsertedID}})
		return
	}
}

func GetBatch() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		batchId := c.Param("id")
		var batch models.Batch
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(batchId)

		err := batchCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&batch)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		c.JSON(http.StatusOK, responses.Response{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": batch}})
		return
	}
}

func UpdateBatch() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		batchId := c.Param("id")
		var batch models.Batch
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(batchId)

		if err := c.BindJSON(&batch); err != nil {
			c.JSON(http.StatusBadRequest, responses.Response{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		if validationErr := validateBatch.Struct(&batch); validationErr != nil {
			c.JSON(http.StatusBadRequest, responses.Response{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}})
			return
		}

		update := bson.M{
			"Vessels":      batch.Vessels,
			"Measurements": batch.Measurements,
			"Volume":       batch.Volume,
		}

		result, err := batchCollection.UpdateOne(ctx, bson.M{"_id": objId}, bson.M{"$set": update})

		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		var updatedBatch models.Batch

		if result.MatchedCount == 1 {
			err := batchCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&updatedBatch)

			if err != nil {
				c.JSON(http.StatusInternalServerError, responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
				return
			}
		}

		c.JSON(http.StatusOK, responses.Response{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": updatedBatch}})
		return

	}
}

func DeleteBatch() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		batchId := c.Param("id")
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(batchId)

		result, err := batchCollection.DeleteOne(ctx, bson.M{"_id": objId})

		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		if result.DeletedCount < 1 {
			c.JSON(http.StatusNotFound, responses.Response{Status: http.StatusNotFound, Message: "error", Data: map[string]interface{}{"data": "batch not found"}})
			return
		}

		c.JSON(http.StatusOK,
			responses.Response{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": "batch deleted"}},
		)
		return
	}
}
