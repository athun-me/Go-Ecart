package controls

import (
	"net/http"
	"strconv"

	"github.com/athunlal/auth"
	"github.com/athunlal/config"

	"github.com/athunlal/models"
	"github.com/gin-gonic/gin"

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

//----------User signup--------------------------------------->

func UserSignUP(c *gin.Context) {
	var Data checkUserData

	if c.Bind(&Data) != nil {
		c.JSON(400, gin.H{
			"error": "Data binding error",
		})
		return
	}

	var temp_user models.User
	hash, err := bcrypt.GenerateFromPassword([]byte(Data.Password), 10)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Status": "False",
			"Error":  "Hashing password error",
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
			c.JSON(500, gin.H{
				"Status": "False",
				"Error":  "User data creating error",
			})
		} else {
			db.Model(&user).Where("email LIKE ?", user.Email).Update("otp", otp)

			c.JSON(202, gin.H{
				"message": "Go to /signup/otpvalidate",
			})
		}

	} else {
		c.JSON(409, gin.H{
			"Error": "User already Exist",
		})
		return
	}
}

//------------User login------------------------>

func UesrLogin(c *gin.Context) {
	type userData struct {
		Email    string
		Password string
	}

	var user userData
	if c.Bind(&user) != nil {
		c.JSON(400, gin.H{
			"error": "Data binding error",
		})
		return
	}
	var checkUser models.User
	db := config.DBconnect()
	result := db.First(&checkUser, "email LIKE ?", user.Email)

	if checkUser.Isblocked == true {
		c.JSON(401, gin.H{
			"Status":  " Authorization",
			"Message": "User blocked by admin",
		})
		return
	}

	if result.Error != nil {
		c.JSON(404, gin.H{
			"Status":  "false",
			"Message": "User not exit",
		})
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(checkUser.Password), []byte(user.Password))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Status": "false",
			"error":  "Password is incorrect",
		})
		return
	}

	//----------------Generating a JWT-tokent-------------------//

	str := strconv.Itoa(int(checkUser.ID))
	tokenString := auth.TokenGeneration(str)
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("UserAutherization", tokenString, 3600*24*30, "", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"message": "User login successfully",
	})

}

//-------------User logout---------------------->

func UserSignout(c *gin.Context) {
	c.SetCookie("UserAutherization", "", -1, "", "", false, false)
	c.JSON(200, gin.H{
		"Message": "User Successfully  Log Out",
	})
}

//-------Validate----------------------->
func Validate(c *gin.Context) {
	c.Get("user")
	c.JSON(200, gin.H{
		"message": "User login successfully",
	})
}
