package controls

import (
	"strconv"

	"github.com/athunlal/config"
	"github.com/athunlal/models"
	"github.com/gin-gonic/gin"
)

type ProfileData struct {
	Firstname   string
	Lastname    string
	Email       string
	PhoneNumber int
}

func GetUserProfile(c *gin.Context) {
	id := c.Param("id")
	var user_data ProfileData
	db := config.DBconnect()
	result := db.Raw("SELECT firstname,lastname,email FROM users WHERE id =?", id).Scan(&user_data)
	if result.Error != nil {
		c.JSON(404, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"Profile Details": user_data,
	})

}
func EditUserProfileByadmin(c *gin.Context) {
	uid := c.Param("id")
	id, err := strconv.Atoi(uid)
	if err != nil {
		c.JSON(400, gin.H{
			"Error": "Error in string conversion",
		})
	}
	var userdata models.User
	if c.Bind(&userdata) != nil {
		c.JSON(400, gin.H{
			"Error": "Unable to Bind JSON data",
		})
		return
	}
	userdata.ID = uint(id)
	DB := config.DBconnect()
	result := DB.Model(&userdata).Updates(models.User{
		Firstname: userdata.Firstname,
		Lastname:  userdata.Lastname,
	})
	if result.Error != nil {
		c.JSON(404, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"Message": "Profile Updated Successfully",
	})

}

func AdminProfile(c *gin.Context) {
	id := c.Param("id")
	var adminData models.Admin
	DB := config.DBconnect()
	result := DB.Raw("SELECT firstname,lastname,email,phone_number FROM admins WHERE id = ?", id).Scan(&adminData)

	if result.Error != nil {
		c.JSON(404, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"Admin Details": adminData,
	})
}
