package controls

import (
	"fmt"
	"time"

	"os"
	"strconv"

	"github.com/athunlal/config"
	"github.com/athunlal/models"
	"github.com/gin-gonic/gin"
	"github.com/razorpay/razorpay-go"
)

func DeleteCartItems(c *gin.Context) {
	id, err := strconv.Atoi(c.GetString("userid"))
	if err != nil {
		c.JSON(400, gin.H{
			"Error": "Error in string conversion",
		})
		return
	}
	// var cartData models.Cart
	db := config.DB
	// result := db.Where("userid = ?", id).Delete(&cartData)
	result := db.Exec("delete from carts where userid = ?", id)

	if result.Error != nil {
		c.JSON(400, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}
}

func CashOnDelivery(c *gin.Context) {
	//fetching user id from token
	id, err := strconv.Atoi(c.GetString("userid"))
	if err != nil {
		c.JSON(400, gin.H{
			"Error": "Error in string conversion",
		})
		return
	}

	var cartData models.Cart

	db := config.DB

	//fetching the data from the table carts by id
	result := db.First(&cartData, "userid = ?", id)
	if result.Error != nil {
		c.JSON(404, gin.H{
			"Message": "Cart is empty",
		})
		return
	}
	//fetching the total amount from the table carts
	var total_amount float64
	result = db.Table("carts").Where("userid = ?", id).Select("SUM(total_price)").Scan(&total_amount)
	if result.Error != nil {
		c.JSON(400, gin.H{
			"Error": "Error fetching the total amount from the table carts",
		})
		return
	}
	todaysDate := time.Now()
	paymentData := models.Payment{
		PaymentMethod: "COD",
		Totalamount:   uint(total_amount),
		Date:          todaysDate,
		Status:        "pending",
		UserId:        uint(id),
	}
	result = db.Create(&paymentData)
	if result.Error != nil {
		c.JSON(400, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}

	var addressData models.Address
	result = db.First(&addressData, "userid = ?", id)
	if result.Error != nil {
		c.JSON(400, gin.H{
			"Error": "address not exist",
		})
		return
	}

	oderData := models.Oder_item{
		UserIdNo:    uint(id),
		TotalAmount: uint(total_amount),
		PaymentId:   paymentData.PaymentId,
		AddId:       addressData.Addressid,
		OrderStatus: "pending",
	}

	result = db.Create(&oderData)
	if result.Error != nil {
		c.JSON(400, gin.H{
			"Error": "Creating the table order",
		})
		return
	}

	c.JSON(200, gin.H{
		"Message": "Payment Method COD",
		"Status":  "True",
	})

	OderDetails(c)
	DeleteCartItems(c)

}

//Online payment using Razorpay
func Razorpay(c *gin.Context) {

	id, err := strconv.Atoi(c.GetString("userid"))
	if err != nil {
		c.JSON(400, gin.H{
			"Error": "Error in string conversion",
		})
	}

	db := config.DB

	//fetching the user data
	var userdata models.User
	result := db.Find(&userdata, "id = ?", id)
	if result.Error != nil {
		c.JSON(404, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}

	//fetching the tatal price from the table carts
	var amount uint
	row := db.Table("carts").Where("userid = ?", id).Select("SUM(total_price)").Row()
	err = row.Scan(&amount)

	if err != nil {
		c.JSON(400, gin.H{
			"Error": err.Error(),
		})
	}

	//Sending the payment details to Razorpay
	client := razorpay.NewClient(os.Getenv("RAZORPAY_KEY_ID"), os.Getenv("RAZORPAY_SECRET"))
	data := map[string]interface{}{
		"amount":   amount * 100,
		"currency": "INR",
		"receipt":  "some_receipt_id",
	}

	//Creating the payment details to client order
	body, err := client.Order.Create(data, nil)
	if err != nil {
		c.JSON(400, gin.H{
			"Error": err,
		})
		return
	}

	//To rendering the html page with user&payment details
	value := body["id"]

	c.HTML(200, "app.html", gin.H{
		"userid":     userdata.ID,
		"totalprice": amount,
		"paymentid":  value,
	})
}

//when the Razorpay payment is completed this funcion will work
func RazorpaySuccess(c *gin.Context) {

	userID, err := strconv.Atoi(c.Query("user_id"))
	if err != nil {
		c.JSON(400, gin.H{
			"Error": "Error in string conversion",
		})
	}

	db := config.DB

	//fetching the payment details from Razorpay
	orderid := c.Query("order_id")
	paymentid := c.Query("payment_id")
	signature := c.Query("signature")
	totalamount := c.Query("total")

	//Creating table razorpay  using the data from Razorpay
	Rpay := models.RazorPay{
		UserID:          uint(userID),
		RazorPaymentId:  paymentid,
		Signature:       signature,
		RazorPayOrderID: orderid,
		AmountPaid:      totalamount,
	}
	result := db.Create(&Rpay)
	if result.Error != nil {
		c.JSON(400, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}

	todyDate := time.Now()
	method := "Razor Pay"
	status := "pending"

	//converting to string total amount veriable
	totalprice, err := strconv.Atoi(totalamount)
	if err != nil {
		c.JSON(400, gin.H{
			"Error": "Error in string conversion",
		})
	}

	//Creating payment table
	paymentdata := models.Payment{
		UserId:        uint(userID),
		PaymentMethod: method,
		Status:        status,
		Date:          todyDate,
		Totalamount:   uint(totalprice),
	}
	result1 := db.Create(&paymentdata)
	if result1.Error != nil {
		c.JSON(400, gin.H{
			"Error": result1.Error.Error(),
		})
		return
	}

	var addressData models.Address
	result = db.First(&addressData, "userid = ?", userID)
	if result.Error != nil {
		c.JSON(400, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}
	pid := paymentdata.PaymentId

	oderData := models.Oder_item{
		UserIdNo:    uint(userID),
		TotalAmount: uint(totalprice),
		PaymentId:   pid,
		AddId:       addressData.Addressid,
		OrderStatus: "pending",
	}

	result = db.Create(&oderData)
	if result.Error != nil {
		c.JSON(400, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"status":    true,
		"paymentid": pid,
	})
	OderDetails(c)
	DeleteCartItems(c)
}

//When the payment is successfull this function will work
func Success(c *gin.Context) {

	pid, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		c.JSON(400, gin.H{
			"Error": "Error in string conversion",
		})
	}

	c.HTML(200, "success.html", gin.H{
		"paymentid": pid,
	})

}

//Wallet payment
func WalletPay(c *gin.Context) {
	//fetching user id from token
	id, err := strconv.Atoi(c.GetString("userid"))
	if err != nil {
		c.JSON(400, gin.H{
			"Error": "Error in string conversion",
		})
	}

	var cartData models.Cart
	var wallet models.Wallet

	db := config.DB

	result := db.Where("user_id = ?", id).First(&wallet)
	if result.Error != nil {
		c.JSON(400, gin.H{
			"Error":    result.Error.Error(),
			"Message ": "Wallet does not exist",
		})
		return
	}

	//fetching the data from the table carts by id
	result = db.First(&cartData, "userid = ?", id)
	if result.Error != nil {
		c.JSON(404, gin.H{
			"Message": "Cart is empty",
		})
		return
	}

	//fetching the total amount from the table carts
	var totalAmount float64
	result = db.Table("carts").Where("userid = ?", id).Select("SUM(total_price)").Scan(&totalAmount)
	if result.Error != nil {
		c.JSON(400, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}

	if wallet.Amount < totalAmount {
		c.JSON(400, gin.H{
			"Error": "Insufficient balance",
		})
		return
	}

	todaysDate := time.Now()
	paymentData := models.Payment{
		PaymentMethod: "Wallet",
		Totalamount:   uint(totalAmount),
		Date:          todaysDate,
		Status:        "pending",
		UserId:        uint(id),
	}
	result = db.Create(&paymentData)
	if result.Error != nil {
		c.JSON(400, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}

	var addressData models.Address
	result = db.First(&addressData, "userid = ?", id)
	if result.Error != nil {
		c.JSON(400, gin.H{
			"Error":   result.Error.Error(),
			"Message": "Address not added",
		})
		return
	}

	oderData := models.Oder_item{
		UserIdNo:    uint(id),
		TotalAmount: uint(totalAmount),
		PaymentId:   paymentData.PaymentId,
		AddId:       addressData.Addressid,
		OrderStatus: "pending",
	}

	result = db.Create(&oderData)
	if result.Error != nil {
		c.JSON(400, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}

	updateAmount := wallet.Amount - totalAmount
	result = db.Model(&wallet).Where("user_id = ?", id).Update("amount", updateAmount)
	if result.Error != nil {
		c.JSON(400, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}

	wHistory := models.WalletHistory{
		UserId:         uint(id),
		Amount:         totalAmount,
		TransctionType: "Debited",
		Date:           todaysDate,
	}
	result = db.Create(&wHistory)
	if result.Error != nil {
		c.JSON(400, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"Message": "Payment Method Wallet",
		"Status":  "True",
	})

	OderDetails(c)
	DeleteCartItems(c)
}

func ShowWallet(c *gin.Context) {
	userid, err := strconv.Atoi(c.GetString("userid"))
	if err != nil {
		c.JSON(400, gin.H{
			"Error": err.Error(),
		})
		return
	}

	fmt.Println("user id : ", userid)
	var wallet models.Wallet

	db := config.DB
	result := db.First(&wallet).Where("user_id = ?", userid)
	if result.Error != nil {
		c.JSON(400, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"Balance Amount": wallet.Amount,
	})
}
