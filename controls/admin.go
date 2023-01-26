package controls

import (
	"net/http"
	"os"
	"strconv"

	"time"

	"github.com/athunlal/auth"
	"github.com/athunlal/config"
	"github.com/athunlal/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type checkAdminData struct {
	Firstname   string
	Lastname    string
	Email       string
	Password    string
	PhoneNumber int
}

func ValidateAdmin(c *gin.Context) {
	c.Get("admin")

	c.JSON(http.StatusOK, gin.H{
		"message": "Admin login successfully",
	})
}

func AdminSignup(c *gin.Context) {
	var Data checkAdminData
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

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": temp_user.Email,
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

func AdminSignout(c *gin.Context) {
	c.SetCookie("AdminAutherization", "", -1, "", "", false, false)
	c.JSON(200, gin.H{
		"Message": "Admin Successfully Signed Out",
	})
}

func AdminLogin(c *gin.Context) {

	type AdminData struct {
		Email    string
		Password string
	}

	var admin AdminData
	if c.Bind(&admin) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Bad request",
		})
		return
	}
	var checkAdmin models.Admin
	db := config.DBconnect()
	result := db.First(&checkAdmin, "email LIKE ?", admin.Email)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"user": "User NOT found",
		})
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(checkAdmin.Password), []byte(admin.Password))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Password is incorrect",
		})
		return
	}

	//----------------Generating a JWT-tokent-------------------//
	str := strconv.Itoa(int(checkAdmin.ID))
	tokenString := auth.TokenGeneration(str)
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("AdminAutherization", tokenString, 3600*24*30, "", "", false, true)

	c.JSON(http.StatusOK, gin.H{})
}

func ViewAllUser(c *gin.Context) {
	var user []models.User
	db := config.DBconnect()
	result := db.Find(&user)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Bad requst",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"User": user,
	})

}
