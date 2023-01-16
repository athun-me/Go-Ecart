package controls

import (
	"net/http"
	"os"
	"time"

	"github.com/athunlal/config"
	"github.com/athunlal/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"

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

//-------validate----------------------->
func Validate(c *gin.Context) {
	user, _ := c.Get("user")
	c.JSON(http.StatusOK, gin.H{
		"message": user,
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

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.Email,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})
	tokenString, err := token.SignedString([]byte(os.Getenv("SECERET")))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to  create token",
		})
		return
	}
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Autherization", tokenString, 3600*24*30, "", "", false, true)

	c.JSON(http.StatusOK, gin.H{})

}
