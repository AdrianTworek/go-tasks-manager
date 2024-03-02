package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email    string `json:"email" gorm:"unique"`
	Password string `json:"password"`
	Tasks    []Task `json:"tasks" gorm:"foreign:UserID"`
}
