package controls

import (
	"net/http"

	"github.com/athunlal/config"
	"github.com/athunlal/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type Data struct {
	Firstname   string
	Lastname    string
	Email       string
	Password    string
	PhoneNumber int
}

func AdminSignup(c *gin.Context) {
	var Data data
	if c.Bind(&Data) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Bad request",
		})
		return
	}

	var temp_user models.Admin
	hash, err := bcrypt.GenerateFromPassword([]byte(Data.Password), 10)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Bad request hashing password",
		})
		return
	}

	// if Data.Otp != Otp {
	// 	c.JSON(http.StatusBadRequest, gin.H{
	// 		"error": "Enter valid OTP",
	// 	})
	// 	return
	// }

	db := config.DBconnect()

	result := db.First(&temp_user, "email LIKE ?", Data.Email)
	if result.Error != nil {
		user := models.Admin{
			Firstname:   Data.Firstname,
			Lastname:    Data.Lastname,
			Email:       Data.Email,
			Password:    string(hash),
			PhoneNumber: Data.PhoneNumber,
		}
		result2 := db.Create(&user)
		if result2.Error != nil {
			c.JSON(http.StatusOK, gin.H{
				"message": "Bad request",
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"message": "Signup successful",
			})
		}

	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "User already Exist",
		})
		return
	}
}
