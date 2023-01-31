package models

type Product struct {
	Productid   uint   `JSON:"productid" gorm:"primarykey;unique"`
	Productname string `JSON:"productname" gorm:"not null"`
	Description string `JSON:"description" gorm:"not null"`
	Stock       uint   `JSON:"stock" gorm:"not null"`
	Price       uint   `JSON:"price" gorm:"not null"`
	Brand       Brand  `gorm:"ForeignKey:Brandid"`
	Brandid     uint   `JSON:"brandid"`
}

type Brand struct {
	ID        uint   `json:"id" gorm:"primaryKey"  `
	Brandname string `JSON:"brandname" gorm:"not null"`
}
