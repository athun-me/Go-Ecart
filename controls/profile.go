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
	id := c.Query("userId")
	var userData models.User
	db := config.DB
	result := db.First(&userData, id)
	if result.Error != nil {
		c.JSON(404, gin.H{
			"Error":   result.Error.Error(),
			"Message": "User not exist",
			"Status":  "false",
		})
		return
	}
	c.JSON(200, gin.H{
		"First name ":  userData.FirstName,
		"Last Name":    userData.LastName,
		"Email":        userData.Email,
		"Phone number": userData.PhoneNumber,
		"Is Block" : userData.Isblocked,
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
	db := config.DB
	result := db.Model(&userData).Updates(models.User{
		FirstName: userEnterdata.Firstname,
		LastName:  userEnterdata.Lastname,

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
	id := c.Query("adminId")
	var user_data models.Admin
	db := config.DB
	result := db.First(&user_data, id)
	if result.Error != nil {
		c.JSON(404, gin.H{
			"Error":   result.Error.Error(),
			"Message": "Admin does't exist",
			"Status":  "false",
		})
		return
	}
	c.JSON(200, gin.H{

		"First name ":  user_data.Firstname,
		"Last Name":    user_data.Lastname,
		"Email":        user_data.Email,
		"Phone number": user_data.PhoneNumber,
	})
}
