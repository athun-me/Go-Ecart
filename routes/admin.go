package routes

import (
	"github.com/athunlal/controls"
	midilware "github.com/athunlal/midileware"
	"github.com/gin-gonic/gin"
)

func AdminRouts(c *gin.Engine) {
	admin := c.Group("/admin")
	{
		admin.POST("/login", controls.AdminLogin)
		admin.POST("/signup", controls.AdminSignup)

		admin.GET("/viewuser", midilware.AdminAuth, controls.ViewAllUser)
		admin.GET("/logout", midilware.AdminAuth, controls.AdminSignout)
		admin.GET("/adminvalidate", midilware.AdminAuth, controls.ValidateAdmin)
	}
	
}
