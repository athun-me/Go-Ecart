package controls

import (
	"strconv"

	"github.com/athunlal/config"
	"github.com/athunlal/models"
	"github.com/gin-gonic/gin"
)

//>>>>>>>>>> Oder Details <<<<<<<<<<<<<<<<

func OderDetails(c *gin.Context) {
	userId, err := strconv.Atoi(c.GetString("userid"))
	if err != nil {
		c.JSON(400, gin.H{
			"Error": "Error in string conversion",
		})
	}

	var UserAddress models.Address
	var UserPayment models.Payment
	var UserCart []models.Cart

	db := config.DBconnect()
	result := db.Find(&UserAddress, "userid = ? AND defaultadd = true", userId)
	if result.Error != nil {
		c.JSON(400, gin.H{
			"Error": result.Error.Error(),
		})
	}
	result = db.Last(&UserPayment, "user_id = ?", userId)
	if result.Error != nil {
		c.JSON(400, gin.H{
			"Error": result.Error.Error(),
		})
	}
	result = db.Find(&UserCart, "userid = ?", userId)
	if result.Error != nil {
		c.JSON(400, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}
	var oder_item models.Oder_item
	db.Last(&oder_item, "useridno = ?", userId)

	for _, UserCart := range UserCart {
		OderDetails := models.OderDetails{
			Userid:     uint(userId),
			AddressId:  UserAddress.Addressid,
			Paymentid:  UserPayment.PaymentId,
			Product_id: UserCart.Product_id,
			Status:     "pending",
			Quantity:   UserCart.Quantity,
			OderItemid: oder_item.OrderId,
		}

		result = db.Create(&OderDetails)
		if result.Error != nil {
			c.JSON(400, gin.H{
				"Error": result.Error.Error(),
			})
			return
		}
	}
	c.JSON(200, gin.H{
		"Message": "Oder Added succesfully",
	})
}

//>>>>>>>>>> Show oder <<<<<<<<<<<<<<<<<<<
func ShowOder(c *gin.Context) {
	userId, err := strconv.Atoi(c.GetString("userid"))
	if err != nil {
		c.JSON(400, gin.H{
			"Error": "Error in string conversion",
		})
		return
	}
	var userOder []models.OderDetails
	var products []models.Product

	db := config.DBconnect()
	result := db.Find(&userOder, "userid = ?", userId)
	if result.Error != nil {
		c.JSON(400, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}

	for _, order := range userOder {

		result := db.Find(&products, "productid = ? ", order.Product_id)
		if result.Error != nil {
			c.JSON(400, gin.H{
				"Error": result.Error.Error(),
			})
			return
		}
		c.JSON(200, gin.H{
			"Product name ": products[0].Productname,
			"Price":         products[0].Price,
			"Description":   products[0].Description,
			"Quantity":      userOder[0].Quantity,
		})
	}
}

//>>>>>>>>>>>>>>< Cancel Oder <<<<<<<<<<<<<<<<<<<<
func CancelOrder(c *gin.Context) {
	userid, err := strconv.Atoi(c.GetString("userid"))
	if err != nil {
		c.JSON(400, gin.H{
			"Error": "Error in string conversion",
		})
	}

	var oder models.OderDetails
	db := config.DBconnect()
	result := db.Model(&oder).Where("userid = ?", userid).Update("status", "Canceled")
	if result.Error != nil {
		c.JSON(400, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"Massage": "oder canceld",
	})
}

//>>>>>>>>>>>>>>< Retrun Oder <<<<<<<<<<<<<<<<<<<
func ReturnOrderByUser(c *gin.Context) {
	userid, err := strconv.Atoi(c.GetString("userid"))
	if err != nil {
		c.JSON(400, gin.H{
			"Error": "Error in string conversion",
		})
		return
	}

	orderId, err := strconv.Atoi(c.Query("orderid"))
	if err != nil {
		c.JSON(400, gin.H{
			"Error": "Error in string conversion",
		})
		return
	}

	var oder models.OderDetails
	var oderItem models.Oder_item
	db := config.DBconnect()
	result := db.Model(&oder).Where("userid = ? AND oder_itemid = ?", userid, orderId).Update("status", "Return product")
	if result.Error != nil {
		c.JSON(400, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}

	result = db.Model(&oderItem).Where("useridno = ? AND order_id = ?", userid, orderId).Update("orderstatus", "Return product")
	if result.Error != nil {
		c.JSON(400, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"Massage": "Product Return",
	})
}

//>>>>>>>>>>>>> Return acsept <<<<<<<<<<<<<<<<<
func ReturnAcsept(c *gin.Context) {

}
