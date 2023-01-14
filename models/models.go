package models

import "gorm.io/gorm"

// type User struct {
// 	gorm.Model
// 	ID          uint   `json:"id" gorm:"primaryKey;unique"  `
// 	Firstname   string `json:"first_name"  gorm:"not null" validate:"required,min=2,max=50"`
// 	Lastname    string `json:"last_name"    gorm:"not null"    validate:"required,min=1,max=50"`
// 	Email       string `json:"email"   gorm:"not null;unique"  validate:"email,required"`
// 	Password    string `json:"password" gorm:"not null"  validate:"required"`
// 	PhoneNumber string `json:"phone"   gorm:"not null;unique" validate:"required"`
// 	Otp         string
// }

type User struct {
	gorm.Model
	ID          uint
	Firstname   string
	Lastname    string
	Email       string
	Password    string
	PhoneNumber string
	Otp         string
}

type Admin struct {
	gorm.Model
	Firstname   string
	Lastname    string
	Email       string
	Password    string
	PhoneNumber string
}
