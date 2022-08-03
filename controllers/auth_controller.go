package controllers

import (
	"aging-api/api"
	"aging-api/auth"
	"aging-api/models"
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
			api.Respond(
				c,
				http.StatusBadRequest,
				"Bad Request",
				err.Error(),
			)
			return
		}

		if err := userCollection.FindOne(ctx, bson.M{"email": requestBody.Email}).Decode(&fetchedUser); err != nil {
			api.Respond(
				c,
				http.StatusNotFound,
				"Account not found",
				err.Error(),
			)
			return
		}

		if authorized := auth.CheckPasswordHash(requestBody.Password, fetchedUser.Password); authorized == false {
			api.Respond(
				c,
				http.StatusUnauthorized,
				"Not authorized",
				"Incorrect password",
			)
			return
		}

		// create JWT token
		token, err := auth.CreateJWT(fetchedUser.Id.Hex())

		if err != nil {
			api.Respond(
				c,
				http.StatusInternalServerError,
				"error",
				err.Error(),
			)
			return
		}
		api.Respond(
			c, http.StatusOK, "success", token,
		)
		return
	}
}
