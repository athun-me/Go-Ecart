package controls

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/athunlal/config"

	"github.com/athunlal/models"
	"github.com/gin-gonic/gin"
)

func Addaddress(c *gin.Context) {
	id := c.Param("id")
	uid, err := strconv.Atoi(id)

	if err != nil {
		c.JSON(501, gin.H{
			"Success": "false",
			"Error":   "Error in string conversion",
		})
	}
	var addressData models.Address

	if c.Bind(&addressData) != nil {
		c.JSON(501, gin.H{
			"Success": "false",
			"Error":   "Error in Binding the JSON",
		})
	}

	addressData.Userid = uint(uid)
	db := config.DBconnect()
	result := db.Create(&addressData)
	if result.Error != nil {
		c.JSON(500, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"Message": "Address added succesfully",
	})
}

func ShowAddress(c *gin.Context) {
	id := c.Param("id")

	var userAddres models.Address
	db := config.DBconnect()
	result := db.First(&userAddres, id)
	if result.Error != nil {
		c.JSON(404, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"User Address": userAddres,
	})
}

func EditUserAddress(c *gin.Context) {
	id := c.Param("id")

	var userAddress models.Address
	if c.Bind(&userAddress) != nil {
		c.JSON(404, gin.H{
			"Error": "Error in binding JSON data",
		})
		return
	}
	str, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(500, gin.H{
			"Error": err,
		})
		return
	}
	userAddress.Userid = uint(str)
	db := config.DBconnect()
	fmt.Println(id)
	result := db.Model(userAddress).Updates(models.Address{
		Name:     userAddress.Name,
		Phoneno:  userAddress.Phoneno,
		Houseno:  userAddress.Houseno,
		Area:     userAddress.Area,
		Landmark: userAddress.Landmark,
		City:     userAddress.City,
		Pincode:  userAddress.Pincode,
		District: userAddress.District,
		State:    userAddress.State,
		Country:  userAddress.Country,
	})

	if result.Error != nil {
		c.JSON(404, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"Message": "Successfully Updated the Address",
		"Updated data" : userAddress,
	})

}
