package models

type User struct {
	Model
	Email    string `json:"email,omitempty" gorm:"unique"`
	Password string `json:"password,omitempty"`
	Tasks    []Task `json:"tasks,omitempty" gorm:"foreign:UserID"`
}
