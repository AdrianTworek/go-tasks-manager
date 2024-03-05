package middleware

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/AdrianTworek/go-tasks-manager/initializers"
	"github.com/AdrianTworek/go-tasks-manager/models"
	"github.com/AdrianTworek/go-tasks-manager/types"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func Authenticate(c *gin.Context) {
	tokenString := c.Request.Header.Get("Authorization")

	if tokenString == "" {
		c.AbortWithStatus(http.StatusUnauthorized)
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_ACCESS_TOKEN_SECRET")), nil
	})

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Could not parse access token",
		})
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		parsedSub, err := uuid.Parse(claims["sub"].(string))
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		var user models.User
		initializers.DB.First(&user, parsedSub)

		if user.ID == uuid.Nil {
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		sanitizedUser := types.SanitizedUser{
			Model: user.Model,
			Email: user.Email,
		}

		c.Set("user", sanitizedUser)
		c.Next()
	} else {
		c.AbortWithStatus(http.StatusUnauthorized)
	}
}
