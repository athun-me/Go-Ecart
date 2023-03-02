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
	var user []models.User
	db := config.DB
	result := db.Limit(3).Find(&user)
	if result.Error != nil {
		c.JSON(500, gin.H{
			"Status":  "False",
			"Message": "Could not find the users",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"User data": user,
	})
}

//<<<<<<<<<<<<<<-Admin search users->>>>>>>>>>>>>>>>>>>>>>>>>

func AdminSearchUser(c *gin.Context) {

	//Get id from URL
	id := c.Param("id")

	var user []models.User
	db := config.DB
	result := db.First(&user, id)

	if result.Error != nil {
		c.JSON(404, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"User data": user,
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
			"Error": result.Error.Error(),
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
