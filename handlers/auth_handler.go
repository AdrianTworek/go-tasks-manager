package handlers

import (
	"net/http"
	"os"
	"time"

	"github.com/AdrianTworek/go-tasks-manager/initializers"
	"github.com/AdrianTworek/go-tasks-manager/models"
	"github.com/AdrianTworek/go-tasks-manager/types"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func HandleLogin(c *gin.Context) {
	var body types.Login

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var user models.User
	initializers.DB.First(&user, "email = ?", body.Email)

	if user.ID == uuid.Nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid credentials",
		})
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid credentials",
		})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	accessToken, err := token.SignedString([]byte(os.Getenv("JWT_ACCESS_TOKEN_SECRET")))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to sign access token",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"accessToken": accessToken})
}
