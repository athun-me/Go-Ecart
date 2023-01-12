package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Firstname   string `json:"name" validate:"required,min=3,max=12"`
	Lastname    string
	Email       string `json:"email" validate:"required, email"`
	Password    string `validate:"requierd, gte=6"`
	PhoneNumber int
}
