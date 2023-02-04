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

//>>>>>>>>>> Get user profile <<<<<<<<<<<<<<<<<<<<<<<<<
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

//>>>>>>>>>>>>>>> Edit user profile <<<<<<<<<<<<<<<<<<
func EditUserProfileByadmin(c *gin.Context) {
	uid := c.Param("id")
	id, err := strconv.Atoi(uid)
	if err != nil {
		c.JSON(400, gin.H{
			"Error": "Error in string conversion",
		})
	}
	var userEnterdata ProfileData
	var userData models.User
	if c.Bind(&userEnterdata) != nil {
		c.JSON(400, gin.H{
			"Error": "Unable to Bind JSON data",
		})
		return
	}
	userData.ID = uint(id)
	db := config.DBconnect()
	result := db.Model(&userData).Updates(models.User{
		Firstname: userEnterdata.Firstname,
		Lastname:  userEnterdata.Lastname,

		PhoneNumber: userEnterdata.PhoneNumber,
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

//>>>>>>>>>>> Admin profile <<<<<<<<<<<<<<<<<<<<<<<
func AdminProfile(c *gin.Context) {
	id := c.Param("id")

	var user_data ProfileData
	db := config.DBconnect()
	result := db.Raw("SELECT firstname,lastname,email,phone_number FROM admins WHERE id =?", id).Scan(&user_data)
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
