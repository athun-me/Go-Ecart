package routes

import (
	"github.com/athunlal/controls"
	midilware "github.com/athunlal/midileware"
	"github.com/gin-gonic/gin"
)

func UserRouts(c *gin.Engine) {
	User := c.Group("/user")
	{
		User.POST("/login", controls.UesrLogin)
		User.POST("/signup", controls.UserSignUP)
		User.POST("/signup/otpvalidate", controls.OtpValidation)

		User.GET("/logout", controls.UserSignout)
		User.GET("/validate", midilware.UserAuth, controls.Validate)

	}
}
