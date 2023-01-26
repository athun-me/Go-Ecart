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
		admin.GET("/adminvalidate", midilware.AdminAuth, controls.ValidateAdmin)
	}
}
