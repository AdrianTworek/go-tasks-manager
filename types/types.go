package types

import (
	"time"

	"github.com/AdrianTworek/go-tasks-manager/models"
	"github.com/google/uuid"
)

type Register struct {
	Email    string `json:"email" binding:"required,email,max=255"`
	Password string `json:"password" binding:"required,min=8,max=32"`
}

type Login struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type SanitizedUser struct {
	models.Model
	Email string
}

type CreateTask struct {
	Title       string `json:"title" binding:"required,max=100"`
	Description string `json:"description" binding:"required,max=1000"`
}

type TaskResponse struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type UpdateTask struct {
	Title       string            `json:"title" binding:"required,max=100"`
	Description string            `json:"description" binding:"required,max=1000"`
	Status      models.TaskStatus `json:"status" binding:"required"`
}
