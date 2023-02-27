package controls

import (
	"fmt"
	"strconv"
	"time"

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

	db := config.DB
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
	db.Last(&oder_item, "user_id_no = ?", userId)

	for _, UserCart := range UserCart {
		OderDetails := models.OderDetails{
			Userid:     uint(userId),
			AddressId:  UserAddress.Addressid,
			PaymentId:  UserPayment.PaymentId,
			ProductId:  UserCart.ProductId,
			Status:     "pending",
			Quantity:   UserCart.Quantity,
			OderItemId: oder_item.OrderId,
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

	db := config.DB
	result := db.Find(&userOder, "userid = ?", userId)
	if result.Error != nil {
		c.JSON(400, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}

	for _, order := range userOder {

		result := db.Find(&products, "product_id = ? ", order.ProductId)
		if result.Error != nil {
			c.JSON(400, gin.H{
				"Error": result.Error.Error(),
			})
			return
		}
		c.JSON(200, gin.H{
			"Product name ": products[0].ProductName,
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

	orderItmeId := c.Query("order_itemid")

	var orderDetails models.OderDetails
	var orderItem models.Oder_item
	var wallet models.Wallet

	db := config.DB

	err = db.First(&orderItem, orderItmeId).Error
	if err != nil {
		c.JSON(400, gin.H{
			"Error": "order id does't exist",
		})
		return
	}

	if orderItem.OrderStatus == "Canceled" {
		c.JSON(400, gin.H{
			"Error": "Oder already canceled",
		})
		return
	}

	result := db.Model(&orderDetails).Where("userid = ? AND oder_item_id = ? ", userid, orderItmeId).Update("status", "Canceled")
	if result.Error != nil {
		c.JSON(400, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}

	result = db.Model(&orderItem).Where("order_id = ?", orderItmeId).Update("order_status", "Canceled")
	if result.Error != nil {
		c.JSON(400, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}

	//adding the balance amount into the wallet
	result = db.Where("user_id = ?", userid).First(&wallet)
	if result.Error != nil {
		walletData := models.Wallet{
			UserId: uint(userid),
			Amount: float64(orderItem.TotalAmount),
		}
		result = db.Create(&walletData)
		if result.Error != nil {
			c.JSON(400, gin.H{
				"Error": result.Error.Error(),
			})
			return
		} else {
			c.JSON(200, gin.H{
				"Message": "Amount added into wallet",
			})
		}
	} else {
		totalAmount := wallet.Amount + float64(orderItem.TotalAmount)
		fmt.Println("this is the added amount : ", totalAmount)

		result = db.Model(&wallet).Where("user_id = ?", userid).Update("amount", totalAmount)
		if result.Error != nil {
			c.JSON(400, gin.H{
				"Error": result.Error.Error(),
			})
			return
		} else {
			c.JSON(200, gin.H{
				"Message": "Amount added into wallet",
			})
		}
	}

	wHistory := models.WalletHistory{
		UserId:         uint(userid),
		Amount:         float64(orderItem.TotalAmount),
		TransctionType: "Credit",
		Date:           time.Now(),
	}

	result = db.Create(&wHistory)
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
	db := config.DB
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
