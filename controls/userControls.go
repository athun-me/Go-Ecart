package controls

import (
	"net/http"
	"strconv"

	"github.com/athunlal/config"
	"github.com/athunlal/models"
	"github.com/gin-gonic/gin"
)

//<<<<<<<<<<-View all first ten users->>>>>>>>>>>>>>>>>>>>>

func ViewAllUser(c *gin.Context) {
	type data struct {
		First_name   string
		Last_name    string
		Email        string
		Isblocked    bool
		Phone_number uint
	}

	// var user models.User
	var userData []data
	db := config.DB
	// result := db.Find(&user).Scan(&userData)
	result := db.Table("users").Select("first_name, last_name, email,isblocked,phone_number").Scan(&userData)
	if result.Error != nil {
		c.JSON(404, gin.H{
			"Message": "Could not find the users",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"User data": userData,
	})
}

//<<<<<<<<<<<<<<-Admin search users->>>>>>>>>>>>>>>>>>>>>>>>>

func AdminSearchUser(c *gin.Context) {

	userid, err := strconv.Atoi(c.Query("userid"))
	if err != nil {
		c.JSON(404, gin.H{
			"Error": err.Error,
		})
		return
	}

	var user models.User
	db := config.DB
	result := db.First(&user, userid)

	if result.Error != nil {
		c.JSON(404, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"First name ":  user.FirstName,
		"Last name ":   user.LastName,
		"Email ":       user.Email,
		"Phone number": user.PhoneNumber,
	})

}

//<<<<<<<<<<-Admin block user->>>>>>>>>>>>>>>>>>>>>>>>>>>>

func AdminBlockUser(c *gin.Context) {
	userid, err := strconv.Atoi(c.Query("userid"))
	if err != nil {
		c.JSON(500, gin.H{
			"Error": "Error occure while converting string",
		})
		return
	}

	var user models.User
	db := config.DB

	result := db.First(&user, userid)
	if result.Error != nil {
		c.JSON(404, gin.H{
			"Message": "User not exist",
		})
		return
	}
	if user.Isblocked == false {
		result := db.Model(&user).Where("id", userid).Update("isblocked", true)
		if result.Error != nil {
			c.JSON(404, gin.H{
				"Error": result.Error.Error(),
			})
			return
		}
		c.JSON(200, gin.H{
			"Message": "User blocked",
		})
	} else {
		result := db.Model(&user).Where("id", userid).Update("isblocked", false)
		if result.Error != nil {
			c.JSON(404, gin.H{
				"Error": result.Error.Error(),
			})
			return
		}
		c.JSON(200, gin.H{
			"Message": "User Unblocked",
		})
	}
}
