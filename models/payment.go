package models

import "time"

type Payment struct {
	Payment_id       uint `JSON:"payment_id" gorm:"primarykey"`
	User             User `gorm:"ForeignKey:User_id"`
	User_id          uint
	PaymentMethod    string `jSON:"payment_method" gorm:"not null"`
	Totalamount      uint   `jSON:"total_amount" gorm:"not null"`
	PaymentPending   bool   `JSON:"status" gorm:"default:false"`
	Paymentcompleted bool   `JSON:"status" gorm:"default:false"`
	PaymentFaild  bool   `JSON:"status" gorm:"default:false"`
}

type OderDetails struct {
	Oderid        uint `JSON:"oderid" gorm:"primarykey"`
	Userid        uint
	User          User `gorm:"ForeignKey:Userid"`
	Address_id    uint
	Address       Address `gorm:"ForeignKey:Address_id"`
	Paymentid     uint
	Payment       Payment `gorm:"ForeignKey:Paymentid"`
	Product_id    uint
	Product       Product `gorm:"ForeignKey:Product_id"`
	Quantity      uint
	OderConfirmed bool `JSON:"OderConfirmed" gorm:"default:false"`
	Shipped       bool `JSON:"Shipped" gorm:"default:false"`
	Deliverd      bool `JSON:"Deliverd" gorm:"default:false"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
