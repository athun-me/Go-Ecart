package controls

import (
	"crypto/rand"
	"fmt"

	"net/http"
	"net/smtp"

	"math/big"
	"strconv"

	"github.com/athunlal/config"
	"github.com/athunlal/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func VerifyOTP(email string) string {
	Otp, err := getRandNum()
	if err != nil {
		panic(err)
	}

	sendMail(email, Otp)
	return Otp
}

func sendMail(email string, otp string) {

	// Sender data.
	from := "golangathunbrototype@gmail.com"
	password := "cpgmhxygkwpirrqi"

	// Receiver email address.
	to := []string{
		email,
	}

	// smtp server configuration.
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	// Message.

	// Authentication.
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Sending email.
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, []byte(otp+" is your otp"))
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Email Sent Successfully!")
}

func getRandNum() (string, error) {
	nBig, e := rand.Int(rand.Reader, big.NewInt(8999))
	if e != nil {
		return "", e
	}
	return strconv.FormatInt(nBig.Int64()+1000, 10), nil
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

type Temp_Data struct {
	Email    string
	Otp      string
	Password string
}

func ForgotPassword(c *gin.Context) {
	var data Temp_Data
	var userData models.User

	if c.Bind(&data) != nil {
		c.JSON(400, gin.H{
			"error": "Data binding error",
		})
		return
	}

	otp := VerifyOTP(data.Email)
	db := config.DBconnect()
	result := db.First(&userData, "email LIKE ?", data.Email)
	if result.Error != nil {
		c.JSON(409, gin.H{
			"Error": "User not exist",
		})
		return
	} else {
		db.Model(&userData).Where("email LIKE ?", data.Email).Update("otp", otp)
		c.JSON(202, gin.H{
			"message": "Go to /forgotpassword",
		})
	}
}

func ForgotPasswordOtpValidation(c *gin.Context) {
	var userEnterData Temp_Data
	var userData models.User
	if c.Bind(&userEnterData) != nil {
		c.JSON(400, gin.H{
			"error": "Data binding error",
		})
		return
	}
	db := config.DBconnect()
	result := db.First(&userData, "otp LIKE ? AND email LIKE ?", userEnterData.Otp, userEnterData.Email)
	if result.Error != nil {
		c.JSON(404, gin.H{
			"Error": result.Error.Error(),
		})
	} else {
		c.JSON(202, gin.H{
			"message": "Go to /changepassword",
		})
	}
}

func ChangePassword(c *gin.Context) {
	var userEnterData Temp_Data
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
	db := config.DBconnect()
	fmt.Println("This is the password", hash)
	result := db.First(&userData, "email LIKE ?", userEnterData.Email)
	if result.Error != nil {
		c.JSON(409, gin.H{
			"Error": "User not exist",
		})
		return
	} else {
		db.Model(&userData).Where("email LIKE ?", userEnterData.Email).Update("password", hash)
		c.JSON(202, gin.H{
			"message": "Successfully updated password",
		})
	}
}
