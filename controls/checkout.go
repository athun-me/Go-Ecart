package controls

import (
	"strconv"

	"github.com/athunlal/config"
	"github.com/gin-gonic/gin"
)

func CheckOut(c *gin.Context) {
	ViewCart(c)

	id, err := strconv.Atoi(c.GetString("userid"))
	if err != nil {
		c.JSON(400, gin.H{
			"Error": "Error while string conversion",
		})
	}
	type data struct {
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
	var userAddressDatas []data
	db := config.DBconnect()
	result := db.Raw("SELECT name, phoneno, houseno, area, landmark, city, pincode,district, state, country FROM addresses WHERE defaultadd = true AND userid = ?", id).Scan(&userAddressDatas)
	if result.Error != nil {
		c.JSON(404, gin.H{
			"Error": result,
		})
		return
	}
	c.JSON(200, gin.H{
		"Default Address of user": userAddressDatas,
	})
	var totalPrice float64
	result1 := db.Table("carts").Where("userid = ?", id).Select("SUM(totalprice)").Scan(&totalPrice).Error
	if result1 != nil {
		c.JSON(400, gin.H{
			"Error": "Can not fetch total amount",
		})
		return
	}
	c.JSON(200, gin.H{
		"Total Price": totalPrice,
	})

}
