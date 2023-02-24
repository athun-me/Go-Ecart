package models

import (
	"time"
)

type Payment struct {
	PaymentId     uint `JSON:"payment_id" gorm:"primarykey"`
	User          User `gorm:"ForeignKey:UserId"`
	UserId        uint
	PaymentMethod string `jSON:"payment_method" gorm:"not null"`
	Totalamount   uint   `jSON:"total_amount" gorm:"not null"`
	Status        string `jSON:"Status" gorm:"not null"`
	Date          time.Time
}

type OderDetails struct {
	Oderid     uint `JSON:"oderid" gorm:"primarykey"`
	Userid     uint
	User       User `gorm:"ForeignKey:Userid"`
	AddressId  uint
	Address    Address `gorm:"ForeignKey:AddressId"`
	PaymentId  uint
	Payment    Payment `gorm:"ForeignKey:PaymentId"`
	OderItemId uint
	ProductId  uint
	Product    Product `gorm:"ForeignKey:ProductId"`
	Quantity   uint
	Status     string `jSON:"Status" gorm:"not null"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type Oder_item struct {
	OrderId     uint    `JSON:"OrderId" gorm:"primarykey"`
	User        User    `gorm:"ForeignKey:UserIdNo"`
	UserIdNo    uint    `json:"useridno"  gorm:"not null" `
	TotalAmount uint    `json:"TotalAmount"  gorm:"not null" `
	Payment     Payment `gorm:"ForeignKey:PaymentId"`
	PaymentId   uint    `json:"PaymentId"`
	OrderStatus string  `json:"orderstatus" `
	Address     Address `gorm:"ForeignKey:AddId"`
	AddId       uint    `json:"addid"  `
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Coupon struct {
	ID            int
	CouponCode    string
	DiscountPrice float64
	CreatedAt     time.Time
	Expired       time.Time
}

type RazorPay struct {
	UserID          uint   `JSON:"userid"`
	RazorPaymentId  string `JSON:"razorpaymentid" gorm:"primaryKey"`
	RazorPayOrderID string `JSON:"razorpayorderid"`
	Signature       string `JSON:"signature"`
	AmountPaid      string `JSON:"amountpaid"`
}

type Wallet struct {
	Id     uint
	User   User `gorm:"ForeignKey:UserId"`
	UserId uint
	Amount float64
}
