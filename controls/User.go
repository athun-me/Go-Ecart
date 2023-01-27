package controls

import (
	
	"net/http"
	"strconv"

	"github.com/athunlal/auth"
	"github.com/athunlal/config"

	"github.com/athunlal/models"
	"github.com/gin-gonic/gin"

	// "github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type checkUserData struct {
	Firstname   string
	Lastname    string
	Email       string
	Password    string
	PhoneNumber int
	Otp         string
}

//-------Validate----------------------->
func Validate(c *gin.Context) {
	c.Get("user")
	c.JSON(http.StatusOK, gin.H{
		"message": "User login successfully",
	})
}

//----------User signup--------------------------------------->

func UserSignUP(c *gin.Context) {
	var Data checkUserData

	if c.Bind(&Data) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Bad request",
		})
		return
	}

	var temp_user models.User
	hash, err := bcrypt.GenerateFromPassword([]byte(Data.Password), 10)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Bad request hashing password",
		})
		return
	}

	otp := VerifyOTP(Data.Email)

	db := config.DBconnect()
	result := db.First(&temp_user, "email LIKE ?", Data.Email)
	if result.Error != nil {
		user := models.User{

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
			db.Model(&user).Where("email LIKE ?", user.Email).Update("otp", otp)

			c.JSON(202, gin.H{
				"message": "Go to /signup/otpvalidate", //202 success but there still one more process
			})
		}

	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "User already Exist",
		})
		return
	}
}

//-------Otp validtioin------------->

func OtpValidation(c *gin.Context) {
	type User_otp struct {
		Otp   string
		Email string
	}
	var user_otp User_otp
	var userData models.User
	if c.Bind(&user_otp) != nil {
		c.JSON(400, gin.H{
			"Error": "Could not bind the JSON Data",
		})
		return
	}
	db := config.DBconnect()
	result := db.First(&userData, "otp LIKE ? AND email LIKE ?", user_otp.Otp, user_otp.Email)
	if result.Error != nil {
		c.JSON(404, gin.H{
			"Error": result.Error.Error(),
		})
		db.Last(&userData).Delete(&userData)
		c.JSON(422, gin.H{
			"Error":   "Wrong OTP Register Once agian",
			"Message": "Goto /signup/otpvalidate",
		})
		return
	}
	c.JSON(200, gin.H{
		"Message": "New User Successfully Registered",
	})

}

//-------------User logout---------------------->

func UserSignout(c *gin.Context) {
	c.SetCookie("Autherization", "", -1, "", "", false, false)
	c.JSON(200, gin.H{
		"Message": "User Successfully  Log Out",
	})
}

//------------User login------------------------>

func UesrLogin(c *gin.Context) {
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


	if checkUser.Isblocked == true{
		c.JSON(http.StatusBadRequest, gin.H{
			"user": "User blocked by admin",
		})
		return
	}

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

	//----------------Generating a JWT-tokent-------------------//

	str := strconv.Itoa(int(checkUser.ID))
	tokenString := auth.TokenGeneration(str)
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Autherization", tokenString, 3600*24*30, "", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"message": "User login successfully",
	})

}
