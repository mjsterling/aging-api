package controllers

import (
	"aging-api/auth"
	"aging-api/models"
	"aging-api/responses"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var requestBody models.User
		var fetchedUser models.User
		defer cancel()

		if err := c.BindJSON(&requestBody); err != nil {
			c.JSON(http.StatusBadRequest, responses.Response{Status: http.StatusBadRequest, Message: "Bad Request", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		if err := userCollection.FindOne(ctx, bson.M{"email": requestBody.Email}).Decode(&fetchedUser); err != nil {
			c.JSON(http.StatusNotFound, responses.Response{Status: http.StatusInternalServerError, Message: "Account not found", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		if authorized := auth.CheckPasswordHash(requestBody.Password, fetchedUser.Password); authorized == false {
			c.JSON(http.StatusUnauthorized, responses.Response{Status: http.StatusUnauthorized, Message: "error", Data: map[string]interface{}{"data": "Incorrect password"}})
			return
		}

		// create JWT token
		token, err := auth.CreateJWT(fetchedUser.Id.Hex())

		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		c.JSON(http.StatusOK,
			responses.Response{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"token": token}},
		)
		return
	}
}
