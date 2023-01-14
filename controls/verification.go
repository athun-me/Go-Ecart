package controls

import (
	"crypto/rand"
	"fmt"
	"net/http"
	"net/smtp"

	// "strings"

	"math/big"
	"strconv"

	"github.com/gin-gonic/gin"
	// "github.com/gin-gonic/gin"
)

var Otp string

func VerifyOTP() {
	Otp, err := getRandNum()
	if err != nil {
		panic(err)
	}
	sendMail(Otp)

}

func sendMail(otp string) {

	// Sender data.
	from := "golangathunbrototype@gmail.com"
	password := "cpgmhxygkwpirrqi"

	// Receiver email address.
	to := []string{
		"athunlalp@gmail.com",
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

func IsUserValid(c *gin.Context) bool {
	VerifyOTP()
	var otp string

	if c.Bind(&otp) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Bad request",
		})
		return false
	}
	
	var OtpVaild bool

	if otp != Otp {
		OtpVaild = false
	} else {
		OtpVaild = true
	}
	return OtpVaild
}
