package controls

import (
	"fmt"

	"os"
	"strconv"

	"github.com/athunlal/config"
	"github.com/athunlal/models"
	"github.com/gin-gonic/gin"
	"github.com/razorpay/razorpay-go"
)

// //>>>>>>>>>>>>>>>> Payment <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<
// func Payment(c *gin.Context) {
// 	//fetching user id from token
// 	id, err := strconv.Atoi(c.GetString("userid"))
// 	if err != nil {
// 		c.JSON(400, gin.H{
// 			"Error": "Error in string conversion",
// 		})
// 	}
// 	type data struct {
// 		Method string
// 	}
// 	var bindData data
// 	var cartDate models.Cart

// 	//binding the data from posman
// 	if c.Bind(&bindData) != nil {
// 		c.JSON(400, gin.H{
// 			"Error": "Could not bind the JSON data",
// 		})
// 		return
// 	}

// 	db := config.DBconnect()

// 	//fetching the data from the table carts by id
// 	result := db.First(&cartDate, "userid = ?", id)
// 	if result.Error != nil {
// 		c.JSON(400, gin.H{
// 			"Error": result.Error.Error(),
// 		})
// 		return
// 	}
// 	//fetching the total amount from the table carts
// 	var total_amount float64
// 	result = db.Table("carts").Where("userid = ?", id).Select("SUM(totalprice)").Scan(&total_amount)
// 	if result.Error != nil {
// 		c.JSON(400, gin.H{
// 			"Error": result.Error.Error(),
// 		})
// 		return
// 	}

// 	if bindData.Method == "COD" {
// 		if result.Error != nil {
// 			c.JSON(400, gin.H{
// 				"Error": result.Error.Error(),
// 			})
// 			return
// 		}

// 		paymentData := models.Payment{
// 			PaymentMethod: bindData.Method,
// 			Totalamount:   uint(total_amount),
// 			User_id:       uint(id),
// 		}

// 		result = db.Create(&paymentData)
// 		if result.Error != nil {
// 			c.JSON(400, gin.H{
// 				"Error": result.Error.Error(),
// 			})
// 			return
// 		}
// 		c.JSON(200, gin.H{
// 			"Message": "Payment Method COD",
// 			"Status":  "Completed",
// 		})

// 	} else if bindData.Method == "UPI" {
// 		//razor pay code
// 		client := razorpay.NewClient("rzp_test_mCL1zwPhJbeuND", "qUeHjny0jl14sphKqOFpyq9M")

// 		data := map[string]interface{}{
// 			"amount":   cartDate.Totalprice,
// 			"currency": "INR",
// 			"receipt":  "some_receipt_id",
// 		}
// 		body, err := client.Order.Create(data, nil)

// 		if err != nil {
// 			c.JSON(500, gin.H{
// 				"Error": err.Error(),
// 			})
// 			return
// 		} else {
// 			paymentData := models.Payment{
// 				PaymentMethod: bindData.Method,
// 				Totalamount:   uint(total_amount),
// 				User_id:       uint(id),
// 			}

// 			result = db.Create(&paymentData)
// 			if result.Error != nil {
// 				c.JSON(400, gin.H{
// 					"Error": result.Error.Error(),
// 				})
// 				return
// 			}
// 			orderID := body["id"].(string)
// 			amount := body["amount"].(float64)
// 			c.JSON(200, gin.H{
// 				"Order ID": orderID,
// 				"Amount":   amount,
// 				"Message":  "Payment Method UPI",
// 				"Status":   "Completed",
// 			})
// 		}
// 	} else {
// 		c.JSON(400, gin.H{
// 			"Error": "Payment field",
// 		})
// 		return
// 	}
// 	OderDetails(c)
// }

//>>>>>>>>> Cash on delivery <<<<<<<<<<<<<<<<<<<<<<<<<<

func CashOnDelivery(c *gin.Context) {
	//fetching user id from token
	id, err := strconv.Atoi(c.GetString("userid"))
	if err != nil {
		c.JSON(400, gin.H{
			"Error": "Error in string conversion",
		})
	}

	var cartData models.Cart

	db := config.DBconnect()

	//fetching the data from the table carts by id
	result := db.First(&cartData, "userid = ?", id)
	if result.Error != nil {
		c.JSON(400, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}
	//fetching the total amount from the table carts
	var total_amount float64
	result = db.Table("carts").Where("userid = ?", id).Select("SUM(totalprice)").Scan(&total_amount)
	if result.Error != nil {
		c.JSON(400, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}
	paymentData := models.Payment{
		PaymentMethod: "COD",
		Totalamount:   uint(total_amount),
		User_id:       uint(id),
	}
	result = db.Create(&paymentData)
	if result.Error != nil {
		c.JSON(400, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"Message": "Payment Method COD",
		"Status":  "Completed",
	})
}

//>>>>>>>>>>>>>> Razorpay <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<

func Razorpay(c *gin.Context) {
	fmt.Println("------------------the first line------------------")
	id, err := strconv.Atoi(c.GetString("userid"))
	if err != nil {
		c.JSON(400, gin.H{
			"Error": "Error in string conversion",
		})
	}

	db := config.DBconnect()
	var userdata models.User
	result := db.Find(&userdata, "id = ?", id)
	if result.Error != nil {
		c.JSON(404, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}
	var amount uint
	result = db.Raw("SELECT SUM(totalprice) FROM carts WHERE id = ?", id).Scan(&amount)
	fmt.Println("this is the total amount : ", amount)
	if result.Error != nil {
		c.JSON(400, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}
	client := razorpay.NewClient(os.Getenv("RAZORPAY_KEY_ID"), os.Getenv("RAZORPAY_SECRET"))
	data := map[string]interface{}{
		"amount":   amount,
		"currency": "INR",
		"receipt":  "some_receipt_id",
	}
	body, err := client.Order.Create(data, nil)
	if err != nil {
		c.JSON(400, gin.H{
			"Error": err,
		})
		return
	}
	value := body["id"]
	c.HTML(200, "app.html", gin.H{
		"userid":      userdata.ID,
		"totalprice":  amount,
		"total":       amount,
		"paymentid":   value,
		"email":       userdata.Email,
		"phonenumber": userdata.PhoneNumber,
	})
	
}

func RazorpaySuccess(c *gin.Context) {
	userid := c.Query("user_id")
	userID, _ := strconv.Atoi(userid)
	orderid := c.Query("order_id")
	paymentid := c.Query("payment_id")
	signature := c.Query("signature")
	totalamount := c.Query("total")
	Rpay := models.RazorPay{
		UserID:          uint(userID),
		RazorPaymentId:  paymentid,
		Signature:       signature,
		RazorPayOrderID: orderid,
		AmountPaid:      totalamount,
	}

	db := config.DBconnect()
	result := db.Create(&Rpay)
	if result.Error != nil {
		c.JSON(400, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}
	method := "Razor Pay"
	status := "pending"
	totalprice, _ := strconv.Atoi(totalamount)
	id, _ := strconv.Atoi(userid)
	paymentdata := models.Payment{
		User_id:       uint(id),
		PaymentMethod: method,
		Status:        status,
		Razorpayid:    paymentid,
		Totalamount:   uint(totalprice),
	}
	result1 := db.Create(&paymentdata)
	if result1.Error != nil {
		c.JSON(400, gin.H{
			"Error": result1.Error.Error(),
		})
		return
	}
	pid := paymentdata.Payment_id
	c.JSON(200, gin.H{

		"status":    true,
		"paymentid": pid,
	})
}

func Success(c *gin.Context) {
	pid, _ := strconv.Atoi(c.Query("id"))
	cid := c.Query("cid")
	fmt.Printf("Fully success assholes")
	c.HTML(200, "success.html", gin.H{
		"paymentid": pid,
		"couponid":  cid,
	})

}
