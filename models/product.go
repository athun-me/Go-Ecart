package models

import "gorm.io/gorm"

type Catogery struct {
	ID           uint `json:"id" gorm:"primaryKey"`
	CatogeryName string
}

type Product struct {
	ProductId   uint     `JSON:"ProductId" gorm:"primarykey;unique"`
	ProductName string   `JSON:"productname" gorm:"not null"`
	Description string   `JSON:"description" gorm:"not null"`
	Stock       uint     `JSON:"stock" gorm:"not null"`
	Price       uint     `JSON:"price" gorm:"not null"`
	Catogery    Catogery `gorm:"ForeignKey:CatogeryId"`
	CatogeryId  uint     `json:"Catogery_id"`
	Brand       Brand    `gorm:"ForeignKey:BrandId"`
	BrandId     uint     `JSON:"BrandId"`
}

type Brand struct {
	ID        uint   `json:"id" gorm:"primaryKey"  `
	BrandName string `JSON:"BrandName" gorm:"not null"`
}

type Cart struct {
	gorm.Model
	Product    Product `gorm:"ForeignKey:ProductId"`
	ProductId  uint
	Quantity   uint
	Price      uint
	TotalPrice uint
	Userid     uint
	User       User `gorm:"ForeignKey:Userid"`
}

type Image struct {
	ID        uint    `json:"id" gorm:"primaryKey"`
	Product   Product `gorm:"ForeignKey:ProductId"`
	ProductId uint    `json:"Product_id"`
	Image     string  `JSON:"Image" gorm:"not null"`
}

type Wishlist struct {
	ID        uint `json:"id" gorm:"primaryKey"`
	Userid    uint
	User      User    `gorm:"ForeignKey:Userid"`
	Product   Product `gorm:"ForeignKey:ProductId"`
	ProductId uint    `json:"Product_id"`
}
