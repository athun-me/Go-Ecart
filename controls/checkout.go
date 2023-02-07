package controls

import (
	"strconv"
	"time"

	"github.com/athunlal/config"
	"github.com/athunlal/models"
	"github.com/gin-gonic/gin"
)

func CheckOut(c *gin.Context) {
	id, err := strconv.Atoi(c.GetString("userid"))
	if err != nil {
		c.JSON(400, gin.H{
			"Error": "Error while string conversion",
		})
	}
	type data struct {
		Coupon   string
		Name     string
		Phoneno  string
		Houseno  string
		Area     string
		Landmark string
		City     string
		Pincode  string
		District string
		State    string
		Country  string
	}
	var userEnterData data
	var coupon models.Coupon
	var discountPercentage float64
	var discountPrice float64
	if c.Bind(&userEnterData) != nil {
		c.JSON(400, gin.H{
			"Error": "Could not bind the JSON data",
		})
		return
	}
	db := config.DBconnect()

	//checking coupon is existig or not
	var count int64
	result := db.Find(&coupon, "coupon_code = ?", userEnterData.Coupon).Count(&count)
	if result.Error != nil {
		c.JSON(400, gin.H{
			"Error": result.Error.Error(),
		})
	}
	if count == 0 {
		c.JSON(400, gin.H{
			"message": "Coupon not exist",
		})
	} else {
		currentTime := time.Now()
		expiredData := coupon.Expired

		if currentTime.Before(expiredData) {
			c.JSON(200, gin.H{
				"message": "Coupon valide",
			})
			discountPercentage = coupon.DiscountPrice
		} else if currentTime.After(expiredData) {
			c.JSON(400, gin.H{
				"message": "Coupon expired",
			})
		}
	}

	//fetching the cart details from the table carts
	ViewCart(c)

	//fetching the data from table addresses
	result = db.Raw("SELECT name, phoneno, houseno, area, landmark, city, pincode,district, state, country FROM addresses WHERE defaultadd = true AND userid = ?", id).Scan(&userEnterData)
	if result.Error != nil {
		c.JSON(404, gin.H{
			"Error": result,
		})
		return
	}
	c.JSON(200, gin.H{
		"Default Address of user": userEnterData,
	})

	//fetching and calculatin the total amount of the cart products
	var totalPrice float64
	result1 := db.Table("carts").Where("userid = ?", id).Select("SUM(totalprice)").Scan(&totalPrice).Error

	//calculating the discount amount
	discountPrice = (discountPercentage / 100) * totalPrice
	totalPriceAfterDeduct := totalPrice - discountPrice

	if result1 != nil {
		c.JSON(400, gin.H{
			"Error": "Can not fetch total amount",
		})
		return
	}
	c.JSON(200, gin.H{
		"Price":        totalPrice,
		"Discount":     discountPrice,
		"Total amount": totalPriceAfterDeduct,
	})
}
