package routes

import (
	"github.com/athunlal/controls"
	midilware "github.com/athunlal/midileware"
	"github.com/gin-gonic/gin"
)

func AdminRouts(c *gin.Engine) {
	admin := c.Group("/admin")
	{
		//Admin rounts
		admin.POST("/login", controls.AdminLogin)
		admin.POST("/signup", controls.AdminSignup)
		admin.GET("/logout", midilware.AdminAuth, controls.AdminSignout)
		admin.GET("/getadminprofile/:id", midilware.AdminAuth, controls.AdminProfile)
		admin.GET("/adminvalidate", midilware.AdminAuth, controls.ValidateAdmin)

		//specification management routes
		admin.POST("/addcatogeries", midilware.UserAuth, controls.AddCatogeries)
		admin.POST("/addbrand", midilware.AdminAuth, controls.AddBrands)
		admin.PUT("/editbrand/:id", midilware.AdminAuth, controls.EditBrand)
		admin.GET("/viewbrandbyadmin", midilware.AdminAuth, controls.ViewBrand)

		//User management routes
		admin.GET("/viewuser", midilware.AdminAuth, controls.ViewAllUser)
		admin.GET("/searchuser/:id", midilware.AdminAuth, controls.AdminSearchUser)
		admin.GET("/getuserprofile/:id", midilware.AdminAuth, controls.GetUserProfile)
		admin.PUT("/edituserprofile/:id", midilware.AdminAuth, controls.EditUserProfileByadmin)
		admin.PUT("/unblockeusers/:id", midilware.AdminAuth, controls.AdminUnlockUser)
		admin.PUT("/blockusers/:id", midilware.AdminAuth, controls.AdminBlockUser)
		
		//product management
		admin.POST("/addimage", midilware.UserAuth, controls.AddImages)
		admin.POST("/addproduct", midilware.AdminAuth, controls.AddProduct)
		
		//coupon routes
		admin.POST("/coupon", midilware.AdminAuth, controls.AddCoupon)
		admin.POST("/checkcoupon", midilware.AdminAuth, controls.CheckCoupon)
		
	}

}
