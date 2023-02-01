package models

import "gorm.io/gorm"

type Product struct {
	Productid   uint   `JSON:"productid" gorm:"primarykey;unique"`
	Productname string `JSON:"productname" gorm:"not null"`
	Description string `JSON:"description" gorm:"not null"`
	Stock       uint   `JSON:"stock" gorm:"not null"`
	Price       uint   `JSON:"price" gorm:"not null"`
	Brand       string `JSON:"brandname"`
	Brandid     uint   `JSON:"brandid"`
}

type Brand struct {
	ID        uint   `json:"id" gorm:"primaryKey"  `
	Brandname string `JSON:"brandname" gorm:"not null"`
}

type Image struct {
	ID      uint    `json:"id" gorm:"primaryKey"`
	Product Product `gorm:"ForeignKey:Pid"`
	Pid     uint    `json:"pid"`
	Image   string  `JSON:"Image" gorm:"not null"`
}

type Cart struct {
	gorm.Model
	Product    Product `gorm:"ForeignKey:Product_id"`
	Product_id uint
	Quantity   uint
	Price      uint
	Totalprice uint
	Cartid     uint
	User       User `gorm:"ForeignKey:Cartid"`
}
