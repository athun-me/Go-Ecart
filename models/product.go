package models

import "gorm.io/gorm"

type Catogery struct {
	ID           uint `json:"id" gorm:"primaryKey"`
	CatogeryName string
}

type Product struct {
	Productid   uint     `JSON:"productid" gorm:"primarykey;unique"`
	Productname string   `JSON:"productname" gorm:"not null"`
	Description string   `JSON:"description" gorm:"not null"`
	Stock       uint     `JSON:"stock" gorm:"not null"`
	Price       uint     `JSON:"price" gorm:"not null"`
	Catogery    Catogery `gorm:"ForeignKey:Catogery_id"`
	Catogery_id uint     `json:"Catogery_id"`
	Brand       Brand    `gorm:"ForeignKey:Brand_id"`
	Brand_id    uint     `JSON:"brandid"`
}

type Brand struct {
	ID        uint   `json:"id" gorm:"primaryKey"  `
	Brandname string `JSON:"brandname" gorm:"not null"`
}

type Cart struct {
	gorm.Model
	Product    Product `gorm:"ForeignKey:Product_id"`
	Product_id uint
	Quantity   uint
	Price      uint
	Totalprice uint
	Userid     uint
	User       User `gorm:"ForeignKey:Userid"`
}

type Image struct {
	ID         uint    `json:"id" gorm:"primaryKey"`
	Product    Product `gorm:"ForeignKey:Product_id"`
	Product_id uint    `json:"Product_id"`
	Image      string  `JSON:"Image" gorm:"not null"`
}

type Wishlist struct {
	ID         uint `json:"id" gorm:"primaryKey"`
	Userid     uint
	User       User    `gorm:"ForeignKey:Userid"`
	Product    Product `gorm:"ForeignKey:Product_id"`
	Product_id uint    `json:"Product_id"`
}
