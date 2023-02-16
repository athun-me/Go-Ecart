package routes

import (
	"github.com/athunlal/controls"
	"github.com/athunlal/middlereware"
	"github.com/gin-gonic/gin"
)

func AdminRouts(c *gin.Engine) {
	admin := c.Group("/admin")
	{
		//Admin rounts
		admin.POST("/login", controls.AdminLogin)
		admin.POST("/signup", controls.AdminSignup)
		admin.GET("/logout", middlereware.AdminAuth, controls.AdminSignout)
		admin.GET("/getadminprofile/:id", middlereware.AdminAuth, controls.AdminProfile)
		admin.GET("/adminvalidate", middlereware.AdminAuth, controls.ValidateAdmin)

		//specification management routes
		admin.POST("/addcatogeries", middlereware.UserAuth, controls.AddCatogeries)
		admin.POST("/addbrand", middlereware.AdminAuth, controls.AddBrands)
		admin.PUT("/editbrand/:id", middlereware.AdminAuth, controls.EditBrand)
		admin.GET("/viewbrandbyadmin", middlereware.AdminAuth, controls.ViewBrand)

		//User management routes
		admin.GET("/viewuser", middlereware.AdminAuth, controls.ViewAllUser)
		admin.GET("/searchuser/:id", middlereware.AdminAuth, controls.AdminSearchUser)
		admin.GET("/getuserprofile/:id", middlereware.AdminAuth, controls.GetUserProfile)
		admin.PUT("/edituserprofile/:id", middlereware.AdminAuth, controls.EditUserProfileByadmin)
		admin.PUT("/unblockeusers/:id", middlereware.AdminAuth, controls.AdminUnlockUser)
		admin.PUT("/blockusers/:id", middlereware.AdminAuth, controls.AdminBlockUser)

		//product management
		admin.POST("/addimage", middlereware.UserAuth, controls.AddImages)
		admin.POST("/addproduct", middlereware.AdminAuth, controls.AddProduct)

		//coupon routes
		admin.POST("/coupon", middlereware.AdminAuth, controls.AddCoupon)
		admin.POST("/checkcoupon", middlereware.AdminAuth, controls.CheckCoupon)

		//Salse Report
		admin.GET("/salesreport", middlereware.AdminAuth, controls.SalesReport)
		admin.GET("/exel", middlereware.AdminAuth, controls.DownloadExel)
		admin.GET("/salsereportpdf/download", middlereware.AdminAuth, controls.Downloadpdf)

	}

}
