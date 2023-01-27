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
		admin.GET("/searchuser/:id", midilware.AdminAuth, controls.AdminSearchUser)
		admin.GET("/adminvalidate", midilware.AdminAuth, controls.ValidateAdmin)

		admin.PUT("/blockusers/:id", midilware.AdminAuth, controls.AdminBlockUser)
		admin.PUT("/unblockeusers/:id", midilware.AdminAuth, controls.AdminUnlockUser)

	}

}
