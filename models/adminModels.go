package models

type Admin struct {
	ID          uint   `json:"id" gorm:"primaryKey;unique"  `
	Firstname   string `json:"first_name"  gorm:"not null" validate:"required,min=2,max=50"`
	Lastname    string `json:"last_name"    gorm:"not null"    validate:"required,min=1,max=50"`
	Email       string `json:"email"   gorm:"not null;unique"  validate:"email,required"`
	Password    string `json:"password" gorm:"not null"  validate:"required"`
	PhoneNumber int    `json:"phone"   gorm:"not null;unique" validate:"required"`
	IsAdmin     bool   `JSON:"isadmin" gorm:"default:true"`

}
