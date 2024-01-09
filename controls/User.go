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
	Firstname       string
	Lastname        string
	Email           string
	Password        string
	ConfirmPassword string
	PhoneNumber     int
	Otp             string
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
	if Data.Email == "" {
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

	db := config.DB
	result := db.First(&temp_user, "email LIKE ?", Data.Email).Error

	if result != nil {
		user := models.User{

			FirstName:   Data.Firstname,
			LastName:    Data.Lastname,
			Email:       Data.Email,
			Password:    string(hash),
			PhoneNumber: Data.PhoneNumber,
		}

		otp := VerifyOTP(Data.Email)
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
	db := config.DB
	result := db.First(&checkUser, "email LIKE ?", user.Email)
	if result.Error != nil {
		c.JSON(404, gin.H{
			"Status":  "false",
			"Message": result.Error.Error(),
		})
		return
	}

	if checkUser.Isblocked == true {
		c.JSON(401, gin.H{
			"Status":  " Authorization",
			"Message": "User blocked by admin",
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

	//>>>>>>>>>>>>>>>>> Generating a JWT-tokent <<<<<<<<<<<<<<<//

	str := strconv.Itoa(int(checkUser.ID))
	tokenString := auth.TokenGeneration(str)
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("UserAutherization", tokenString, 3600*24*30, "", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"message": "User login successfully",
	})

}

//>>>>>>>>>>>>>>> Validate <<<<<<<<<<<<<<<<<<<<<<<<<<<<
func Validate(c *gin.Context) {
	c.Get("user")
	c.JSON(200, gin.H{
		"message": "User login successfully",
	})
}

//>>>>>>>>>>>> Change password by user <<<<<<<<<<<<<<<<

func UserChangePassword(c *gin.Context) {

	var userEnterData checkUserData
	if c.Bind(&userEnterData) != nil {
		c.JSON(400, gin.H{
			"error": "Data binding error",
		})
		return
	}

	if userEnterData.Password != userEnterData.ConfirmPassword {
		c.JSON(400, gin.H{
			"Message": "Password not match",
		})
		return
	}

	id, err := strconv.Atoi(c.GetString("userid"))

	if err != nil {
		c.JSON(500, gin.H{
			"Error": "Error in string conversion",
		})
		return
	}

	var userData models.User
	db := config.DB
	result := db.First(&userData, "id = ?", id)
	if result.Error != nil {
		c.JSON(409, gin.H{
			"Error": "User not exist",
		})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(userData.Password), []byte(userEnterData.Password))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "false",
			"error":  "Password is incorrect",
		})
		return
	} else {
		c.JSON(200, gin.H{
			"message": "Go to /useraccess/updatepassword",
		})
	}
}

//>>>>>>>>>>>>> updating new user <<<<<<<<<<<<<<<<<<<<<<
func Updatepassword(c *gin.Context) {

	var userEnterData checkUserData
	var userData models.User
	if c.Bind(&userEnterData) != nil {
		c.JSON(400, gin.H{
			"error": "Data binding error",
		})
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(userEnterData.Password), 10)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Status": "False",
			"Error":  "Hashing password error",
		})
		return
	}

	id, err := strconv.Atoi(c.GetString("userid"))
	if err != nil {
		c.JSON(400, gin.H{
			"Error": "Error in string conversion",
		})
		return
	}
	db := config.DB
	result := db.First(&userData, "id = ?", id)
	if result.Error != nil {
		c.JSON(409, gin.H{
			"Error": "User not exist",
		})
		return
	}

	db.Model(&userData).Where("id = ?", id).Update("password", hash)
	c.JSON(202, gin.H{
		"message": "Successfully updated password",
	})
}

//>>>>>>>>>>>>Veiw user profile <<<<<<<<<<<<<<<<<<<<<<<<<<
func ShowUserDetails(c *gin.Context) {
	var userData models.User
	id, err := strconv.Atoi(c.GetString("userid"))
	if err != nil {
		c.JSON(400, gin.H{
			"Error": "Error in string conversion",
		})
		return
	}
	db := config.DB
	result := db.First(&userData, "id = ?", id)
	if result.Error != nil {
		c.JSON(409, gin.H{
			"Error": "User not exist",
		})
		return
	}
	c.JSON(202, gin.H{

		"First name":   userData.FirstName,
		"Last name":    userData.LastName,
		"Emial":        userData.Email,
		"Phone number": userData.PhoneNumber,
	})
}

//>>>>>>>>>>> Edit user profile <<<<<<<<<<<<<<<<<<<<<<<<<<<<<

func EditUserProfilebyUser(c *gin.Context) {
	var userEnterData checkUserData
	if c.Bind(&userEnterData) != nil {
		c.JSON(400, gin.H{
			"error": "Data binding error",
		})
		return
	}

	var userData models.User
	id, err := strconv.Atoi(c.GetString("userid"))
	if err != nil {
		c.JSON(400, gin.H{
			"Error": "Error in string conversion",
		})
		return
	}
	db := config.DB
	result := db.First(&userData, "id = ?", id)
	if result.Error != nil {
		c.JSON(409, gin.H{
			"Error": "User not exist",
		})
		return
	}
	result = db.Model(&userData).Updates(models.User{
		FirstName:   userEnterData.Firstname,
		LastName:    userEnterData.Lastname,
		PhoneNumber: userEnterData.PhoneNumber,
	})

	if result.Error != nil {
		c.JSON(404, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"Message": "Successfully Updated the profile",
		"Updated data": gin.H{
			"First name": userData.FirstName,
			"Last name":  userData.LastName,
			"Phone":      userData.PhoneNumber,
		},
	})
}
