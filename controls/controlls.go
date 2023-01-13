package controls

import (
	"fmt"
	"net/http"

	"github.com/athunlal/config"
	"github.com/athunlal/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func UserSignUP(c *gin.Context) {
	type data struct {
		Firstname   string
		Lastname    string
		Email       string
		Password    string
		PhoneNumber int
	}
	var Data data
	var temp_user models.User
	if c.Bind(&Data) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Bad request",
		})
		return
	}
	fmt.Println(Data.Firstname)
	hash, err := bcrypt.GenerateFromPassword([]byte(Data.Password), 10)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Bad request hashing password",
		})
		return
	}

	db := config.DBconnect()

	result := db.First(&temp_user, "email LIKE ?", Data.Email)
	if result.Error != nil {
		user := models.User{
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

func Login(c *gin.Context) {
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
	var checkUser models.User
	db := config.DBconnect()
	result := db.First(&checkUser, "email LIKE ?", user.Email)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"user": "User NOT found",
		})
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(checkUser.Password), []byte(user.Password))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Password is incorrect",
		})
		return
	}

	c.JSON(http.StatusBadRequest, gin.H{
		"user": checkUser,
	})
}
