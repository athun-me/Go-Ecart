package models

import (
	"time"
)

type Payment struct {
	Payment_id    uint `JSON:"payment_id" gorm:"primarykey"`
	User          User `gorm:"ForeignKey:User_id"`
	User_id       uint
	PaymentMethod string `jSON:"payment_method" gorm:"not null"`
	Totalamount   uint   `jSON:"total_amount" gorm:"not null"`
	Status        string `jSON:"Status" gorm:"not null"`
	// RazorPay      RazorPay  `gorm:"ForeignKey:razorpayid"`
	// Razorpayid    string    `JSON:"razorpayid" gorm:"defualt:null"`
}

type OderDetails struct {
	Oderid     uint `JSON:"oderid" gorm:"primarykey"`
	Userid     uint
	User       User `gorm:"ForeignKey:Userid"`
	Address_id uint
	Address    Address `gorm:"ForeignKey:Address_id"`
	Paymentid  uint
	Payment    Payment `gorm:"ForeignKey:Paymentid"`
	Product_id uint
	Product    Product `gorm:"ForeignKey:Product_id"`
	Quantity   uint
	Status     string `jSON:"Status" gorm:"not null"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
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
