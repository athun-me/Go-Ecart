package controls

import (
	"net/http"

	"github.com/athunlal/config"
	"github.com/athunlal/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Data struct {
	Firstname   string
	Lastname    string
	Email       string
	Password    string
	PhoneNumber string
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
			Model:       gorm.Model{},
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

func AdminLogin(c *gin.Context) {
	// type checkAdminData struct {
	// 	Email    string
	// 	Password string
	// }

	// var user checkAdminData
	// if c.Bind(&user) != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{
	// 		"error": "Bad request",
	// 	})
	// 	return
	// }

	// var adminData models.Admin
	// db := config.DBconnect()
	// result := db.First(&adminData, "email LIKE ?", user.Email)
	// if result != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{
	// 		"error": "User not found",
	// 	})
	// 	return
	// }
	// err := bcrypt.CompareHashAndPassword([]byte(adminData.Password), []byte(user.Password))
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{
	// 		"error": "error password hashing",
	// 	})
	// 	return
	// }
	// c.JSON(http.StatusBadRequest, gin.H{
	// 	"user": adminData,
	// })

	type userData struct {
		Email    string
		Password string
	}

	var user userData
	if c.Bind(&user) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Bad request",
		})
		return
	}
	var checkAdmin models.Admin
	db := config.DBconnect()
	result := db.First(&checkAdmin, "email LIKE ?", user.Email)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"user": "User NOT found",
		})
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(checkAdmin.Password), []byte(user.Password))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Password is incorrect",
		})
		return
	}

	c.JSON(http.StatusBadRequest, gin.H{
		"user": checkAdmin,
	})
}
