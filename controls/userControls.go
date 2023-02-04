package controls

import (
	"fmt"
	"net/http"

	"github.com/athunlal/config"
	"github.com/athunlal/models"
	"github.com/gin-gonic/gin"
)

//<<<<<<<<<<-View all first ten users->>>>>>>>>>>>>>>>>>>>>

func ViewAllUser(c *gin.Context) {
	var user []models.User
	db := config.DBconnect()
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
	db := config.DBconnect()
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
	id := c.Param("id")

	var user models.User
	db := config.DBconnect()
	var count int64

	result := db.Model(user).Where("id = ?", id).Update("isblocked", true).Count(&count)
	if count == 0 {
		c.JSON(500, gin.H{
			"Message": "Could not find the users",
		})
		return
	}
	if result.Error != nil {
		c.JSON(404, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}
	fmt.Println(user.Firstname)
	c.JSON(200, gin.H{
		"Massage": "Blocked",
	})

}

//<<<<<<<<<<-Admin unblock user->>>>>>>>>>>>>>>>>>>>>>>>>>>>>

func AdminUnlockUser(c *gin.Context) {
	id := c.Param("id")

	var user []models.User
	db := config.DBconnect()
	var count int64
	result := db.Model(user).Where("id = ?", id).Update("isblocked", false).Count(&count)
	if count ==0 {
		c.JSON(500, gin.H{
			"Message": "Could not find the users",
		})
		return
	}
	if result.Error != nil {
		c.JSON(404, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"Massage": "Unblocked",
	})

}
