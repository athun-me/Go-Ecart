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

type checkAdminData struct {
	Firstname   string
	Lastname    string
	Email       string
	Password    string
	PhoneNumber int
}

func ValidateAdmin(c *gin.Context) {
	c.Get("admin")

	c.JSON(200, gin.H{
		"message": "Admin login successfully",
	})
}

//----------------Admin signup-------------------

func AdminSignup(c *gin.Context) {
	var Data checkAdminData
	if c.Bind(&Data) != nil {
		c.JSON(400, gin.H{
			"error": "Data binding error",
		})
		return
	}

	var temp_user models.Admin
	hash, err := bcrypt.GenerateFromPassword([]byte(Data.Password), 10)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Hashing password error",
		})
		return
	}

	db := config.DB

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
			c.JSON(500, gin.H{
				"error": "Admin data creating error",
			})
		} else {
			c.JSON(200, gin.H{
				"message": "Signup successful",
			})
		}

	} else {
		c.JSON(409, gin.H{
			"error": "User already Exist",
		})
		return
	}

}

//----------------Admin signout-------------------

func AdminSignout(c *gin.Context) {
	c.SetCookie("AdminAutherization", "", -1, "", "", false, false)
	c.JSON(200, gin.H{
		"Message": "Admin Successfully Signed Out",
	})
}

//----------------Admin Login-------------------

func AdminLogin(c *gin.Context) {

	type AdminData struct {
		Email    string
		Password string
	}

	var admin AdminData

	if c.Bind(&admin) != nil {
		c.JSON(400, gin.H{
			"error": "Login data binding error",
		})
		return
	}

	var checkAdmin models.Admin

	db := config.DB
	result := db.First(&checkAdmin, "email LIKE ?", admin.Email)
	if result.Error != nil {
		c.JSON(404, gin.H{
			"error": "User not found",
		})
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(checkAdmin.Password), []byte(admin.Password))
	if err != nil {
		c.JSON(501, gin.H{
			"error": "Username and password invalid",
		})
		return
	}

	//----------------Generating a JWT-tokent-------------------
	str := strconv.Itoa(int(checkAdmin.ID))
	tokenString := auth.TokenGeneration(str)
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("AdminAutherization", tokenString, 3600*24*30, "", "", false, true)

	c.JSON(200, gin.H{
		"message": "Successfully Login to Admin panel",
	})
}
