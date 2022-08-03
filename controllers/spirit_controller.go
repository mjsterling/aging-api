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

var spiritCollection *mongo.Collection = configs.GetCollection(configs.DB, "spirits")
var validateSpirit = validator.New()

func CreateSpirit() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var spirit models.Spirit
		defer cancel()

		if err := c.BindJSON(&spirit); err != nil {
			c.JSON(http.StatusBadRequest,
				responses.Response{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}},
			)
			return
		}

		if validationErr := validateSpirit.Struct(&spirit); validationErr != nil {
			c.JSON(http.StatusBadRequest, responses.Response{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}})
			return
		}

		newSpirit := models.Spirit{
			Id:         primitive.NewObjectID(),
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
			c.JSON(http.StatusInternalServerError, responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		c.JSON(http.StatusCreated, responses.Response{Status: http.StatusCreated, Message: "success", Data: map[string]interface{}{"data": result.InsertedID}})
		return
	}
}

func GetAllSpirits() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		cur, err := userCollection.Find(ctx, bson.D{})
		results := make([]interface{}, 0)
		for cur.Next(ctx) {
			var result bson.D
			if err := cur.Decode(&result); err != nil {
				c.JSON(http.StatusInternalServerError, responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
				return
			}
			results = append(results, result)
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		c.JSON(http.StatusOK, responses.Response{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"spirits": results}})
		return
	}
}

func GetSpiritById() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		spiritId := c.Param("id")
		var spirit models.Spirit
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(spiritId)

		err := spiritCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&spirit)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		c.JSON(http.StatusOK, responses.Response{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": spirit}})
		return
	}
}

func UpdateSpirit() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		spiritId := c.Param("id")
		var spirit models.Spirit
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(spiritId)

		if err := c.BindJSON(&spirit); err != nil {
			c.JSON(http.StatusBadRequest, responses.Response{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		if validationErr := validateSpirit.Struct(&spirit); validationErr != nil {
			c.JSON(http.StatusBadRequest, responses.Response{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}})
			return
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
			c.JSON(http.StatusInternalServerError, responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		var updatedSpirit models.Spirit

		if result.MatchedCount == 1 {
			err := spiritCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&updatedSpirit)

			if err != nil {
				c.JSON(http.StatusInternalServerError, responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
				return
			}
		}

		c.JSON(http.StatusOK, responses.Response{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": updatedSpirit}})
		return

	}
}

func DeleteSpirit() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		spiritId := c.Param("id")
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(spiritId)

		result, err := spiritCollection.DeleteOne(ctx, bson.M{"_id": objId})

		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		if result.DeletedCount < 1 {
			c.JSON(http.StatusNotFound, responses.Response{Status: http.StatusNotFound, Message: "error", Data: map[string]interface{}{"data": "spirit not found"}})
			return
		}

		c.JSON(http.StatusOK,
			responses.Response{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": "spirit deleted"}},
		)
		return
	}
}
