package auth

import (
	"aging-api/api"
	"net/http"
	"os"
	"regexp"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func CreateJWT(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": userID,
		"exp":    time.Now().Add(time.Hour * 24 * 7).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("SESSION_SECRET")))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func DecodeJwt(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SESSION_SECRET")), nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims["userID"].(string), nil
	}

	return "", nil
}

func Authenticate(c *gin.Context) bool {
	authHeader := c.Request.Header["Authorization"]
	if authHeader == nil {
		api.Respond(c,
			http.StatusUnauthorized,
			"Error: No JWT header present",
			"Error: No JWT header present",
		)
		return false
	}
	r := regexp.MustCompile("Bearer (.+)")
	jwt := r.FindStringSubmatch(authHeader[0])
	_, err := DecodeJwt(jwt[1])
	if err != nil {
		api.Respond(c,
			http.StatusUnauthorized,
			"Error: JWT not valid",
			"Error: JWT not valid",
		)
	}
	return true
}
