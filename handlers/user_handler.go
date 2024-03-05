package handlers

import (
	"net/http"

	"github.com/AdrianTworek/go-tasks-manager/initializers"
	"github.com/AdrianTworek/go-tasks-manager/models"
	"github.com/AdrianTworek/go-tasks-manager/types"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func HandleRegister(c *gin.Context) {
	var body types.Register

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var existingUser models.User
	initializers.DB.First(&existingUser, "email = ?", body.Email)

	if existingUser.ID != uuid.Nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": "User already exists",
		})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), 12)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to hash password",
		})
		return
	}

	user := models.User{Email: body.Email, Password: string(hashedPassword)}
	result := initializers.DB.Create(&user)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create user",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}
