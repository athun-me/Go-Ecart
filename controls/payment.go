package controls

import (
	"strconv"

	"github.com/athunlal/config"
	"github.com/athunlal/models"
	"github.com/gin-gonic/gin"
	"github.com/razorpay/razorpay-go"
)

//>>>>>>>>>>>>>>>> Payment <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<
func Payment(c *gin.Context) {
	//fetching user id from token
	id, err := strconv.Atoi(c.GetString("userid"))
	if err != nil {
		c.JSON(400, gin.H{
			"Error": "Error in string conversion",
		})
	}
	type data struct {
		Method string
	}
	var bindData data
	var cartDate models.Cart

	//binding the data from posman
	if c.Bind(&bindData) != nil {
		c.JSON(400, gin.H{
			"Error": "Could not bind the JSON data",
		})
		return
	}

	db := config.DBconnect()

	//fetching the data from the table carts by id
	result := db.First(&cartDate, "userid = ?", id)
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

	if bindData.Method == "COD" {
		if result.Error != nil {
			c.JSON(400, gin.H{
				"Error": result.Error.Error(),
			})
			return
		}

		paymentData := models.Payment{
			PaymentMethod: bindData.Method,
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

	} else if bindData.Method == "UPI" {
		//razor pay code
		client := razorpay.NewClient("rzp_test_mCL1zwPhJbeuND", "qUeHjny0jl14sphKqOFpyq9M")

		data := map[string]interface{}{
			"amount":   cartDate.Totalprice,
			"currency": "INR",
			"receipt":  "some_receipt_id",
		}
		body, err := client.Order.Create(data, nil)

		if err != nil {
			c.JSON(500, gin.H{
				"Error": err.Error(),
			})
			return
		} else {
			paymentData := models.Payment{
				PaymentMethod: bindData.Method,
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
			orderID := body["id"].(string)
			amount := body["amount"].(float64)
			c.JSON(200, gin.H{
				"Order ID": orderID,
				"Amount":   amount,
				"Message":  "Payment Method UPI",
				"Status":   "Completed",
			})
		}
	} else {
		c.JSON(400, gin.H{
			"Error": "Payment field",
		})
		return
	}
	OderDetails(c)
}
