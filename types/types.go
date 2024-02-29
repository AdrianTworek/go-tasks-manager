package types

import "gorm.io/gorm"

type Register struct {
	Email    string `json:"email" binding:"required,email,max=255"`
	Password string `json:"password" binding:"required,min=8,max=32"`
}

type Login struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type SanitizedUser struct {
	gorm.Model
	Email string
}
