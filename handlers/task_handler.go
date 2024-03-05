package handlers

import (
	"net/http"

	"github.com/AdrianTworek/go-tasks-manager/initializers"
	"github.com/AdrianTworek/go-tasks-manager/models"
	"github.com/AdrianTworek/go-tasks-manager/types"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func HandleCreateTask(c *gin.Context) {
	var body types.CreateTask

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	user, _ := c.Get("user")
	currentUser, ok := user.(types.SanitizedUser)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to convert user",
		})
		return
	}

	task := models.Task{Title: body.Title, Description: body.Description, UserID: currentUser.ID}
	result := initializers.DB.Create(&task)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create task",
		})
		return
	}

	c.JSON(http.StatusCreated, task)
}

func HandleGetUserTasks(c *gin.Context) {
	user, _ := c.Get("user")
	currentUser, ok := user.(types.SanitizedUser)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to convert user",
		})
		return
	}

	var tasks []models.Task
	result := initializers.DB.Select("id, title, description, status, created_at, updated_at").Where("user_id = ?", currentUser.ID).Find(&tasks)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch user tasks",
		})
		return
	}

	var taskResponses []types.TaskResponse

	for _, task := range tasks {
		taskResponses = append(taskResponses, types.TaskResponse{
			ID:          task.ID,
			Title:       task.Title,
			Description: task.Description,
			Status:      string(task.Status),
			CreatedAt:   task.CreatedAt,
			UpdatedAt:   task.UpdatedAt,
		})
	}

	c.JSON(http.StatusOK, taskResponses)
}

func HandleGetTask(c *gin.Context) {
	user, _ := c.Get("user")
	currentUser, ok := user.(types.SanitizedUser)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to convert user",
		})
		return
	}

	taskId, ok := c.Params.Get("taskId")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Missing task id parameter",
		})
		return
	}

	parsedTaskId, err := uuid.Parse(taskId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Task id must be of type uuid",
		})
		return
	}

	var task models.Task
	result := initializers.DB.First(&task, parsedTaskId)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch task",
		})
		return
	}

	if currentUser.ID != task.UserID {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Forbidden",
		})
		return
	}

	c.JSON(http.StatusCreated, task)
}

func HandleUpdateTask(c *gin.Context) {
	user, _ := c.Get("user")
	currentUser, ok := user.(types.SanitizedUser)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to convert user",
		})
		return
	}

	taskId, ok := c.Params.Get("taskId")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Missing task id parameter",
		})
		return
	}

	parsedTaskId, err := uuid.Parse(taskId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Task id must be of type uuid",
		})
		return
	}

	var task models.Task
	result := initializers.DB.First(&task, parsedTaskId)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch task",
		})
		return
	}

	if currentUser.ID != task.UserID {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Forbidden",
		})
		return
	}

	var body types.UpdateTask

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if !(body.Status == models.TaskTodo || body.Status == models.TaskInProgress || body.Status == models.TaskDone) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Wrong task status",
		})
		return
	}

	task.Title = body.Title
	task.Description = body.Description
	task.Status = body.Status

	if err := initializers.DB.Save(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update task"})
		return
	}

	c.JSON(http.StatusOK, task)
}

func HandleDeleteTask(c *gin.Context) {
	user, _ := c.Get("user")
	currentUser, ok := user.(types.SanitizedUser)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to convert user",
		})
		return
	}

	taskId, ok := c.Params.Get("taskId")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Missing task id parameter",
		})
		return
	}

	parsedTaskId, err := uuid.Parse(taskId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Task id must be of type uuid",
		})
		return
	}

	var task models.Task
	result := initializers.DB.First(&task, parsedTaskId)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch task",
		})
		return
	}

	if currentUser.ID != task.UserID {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Forbidden",
		})
		return
	}

	result = initializers.DB.Delete(&models.Task{}, parsedTaskId)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete task",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Task deleted successfully",
	})
}
