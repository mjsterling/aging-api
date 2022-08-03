package api

import (
	"aging-api/responses"

	"github.com/gin-gonic/gin"
)

func Respond(c *gin.Context, Status int, Message string, Data interface{}) {
	c.JSON(Status, responses.Response{Status: Status, Message: Message, Data: map[string]interface{}{"data": Data}})
}
