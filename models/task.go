package models

import "gorm.io/gorm"

type TaskStatus string

const (
	TaskTodo       TaskStatus = "TODO"
	TaskInProgress TaskStatus = "IN_PROGRESS"
	TaskDone       TaskStatus = "DONE"
)

type Task struct {
	gorm.Model
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      TaskStatus `json:"status" gorm:"default:TODO"`
	UserID      uint       `json:"user_id"`
	User        User       `json:"user" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
