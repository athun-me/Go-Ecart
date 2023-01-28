package models

type User struct {
	ID          uint   `json:"id" gorm:"primaryKey;unique"  `
	Firstname   string `json:"first_name"  gorm:"not null" validate:"required,min=2,max=50"`
	Lastname    string `json:"last_name"    gorm:"not null"    validate:"required,min=1,max=50"`
	Email       string `json:"email"   gorm:"not null;unique"  validate:"email,required"`
	Password    string `json:"password" gorm:"not null"  validate:"required"`
	PhoneNumber int    `json:"phone"   gorm:"not null;unique" validate:"required"`
	IsAdmin     bool   `JSON:"isadmin" gorm:"default:false"`
	Otp         string `JSON:"otp"`
	Isblocked   bool   `JSON:"isblocked" gorm:"default:false"`
}

type Address struct {
	Addressid uint   `JSON:"addressid" gorm:"primarykey;unique"`
	Userid    uint   `JSON:"userid" gorm:"foreignKey:UserRefer"`
	Name      string `JSON:"name" gorm:"not null"`
	Phoneno   string `JSON:"phoneno" gorm:"not null"`
	Houseno   string `JSON:"houseno" gorm:"not null"`
	Area      string `JSON:"area" gorm:"not null"`
	Landmark  string `JSON:"landmark" gorm:"not null"`
	City      string `JSON:"city" gorm:"not null"`
	Pincode   string `JSON:"pincode" gorm:"not null"`
	District  string `JSON:"district" gorm:"not null"`
	State     string `JSON:"state" gorm:"not null"`
	Country   string `JSON:"country" gorm:"not null"`
}
