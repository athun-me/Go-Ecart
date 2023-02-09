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
		User.POST("/addaddress", midilware.UserAuth, controls.Addaddress)
		User.POST("/addtocart", midilware.UserAuth, controls.AddToCart)
		User.POST("/changepassword", midilware.UserAuth, controls.ChangePassword)
		User.POST("/userchangepassword", midilware.UserAuth, controls.UserChangePassword)
		User.POST("/checkcoupon", midilware.UserAuth, controls.CheckCoupon)
		User.POST("/applycoupon", midilware.UserAuth, controls.Applycoupon)

		User.PUT("/editaddress/:id", midilware.UserAuth, controls.EditUserAddress)
		User.PUT("/forgotpassword", midilware.UserAuth, controls.ForgotPassword)
		User.PUT("/forgotpasswordotpvalidation", midilware.UserAuth, controls.ForgotPasswordOtpValidation)
		User.PUT("/updatepassword", midilware.UserAuth, controls.Updatepassword)
		User.PUT("/editprofile", midilware.UserAuth, controls.EditUserProfilebyUser)

		
		User.GET("/checkout", midilware.UserAuth, controls.CheckOut)
		User.GET("/wishlist/:id", midilware.UserAuth, controls.Wishlist)
		User.GET("/viewproducts", midilware.UserAuth, controls.ViewProducts)
		User.GET("/viewbrand", midilware.UserAuth, controls.ViewBrand)
		User.GET("/viewcart", midilware.UserAuth, controls.ViewCart)
		User.GET("/viewprofile", midilware.UserAuth, controls.ShowUserDetails)
		User.GET("/searchaddress/:id", midilware.UserAuth, controls.ShowAddress)
		User.GET("/logout", midilware.UserAuth, controls.UserSignout)
		User.GET("/validate", midilware.UserAuth, controls.Validate)
		User.GET("/fileterbycatogery/:id", midilware.UserAuth, controls.FilteringByCatogery)
		User.GET("/search", midilware.UserAuth, controls.SearchProduct)
		User.GET("/showoder", midilware.UserAuth, controls.ShowOder)
		User.GET("/canceloder", midilware.UserAuth, controls.CancelOder)
		User.GET("/return", midilware.UserAuth, controls.ReturnOder)

		User.DELETE("/deletecart/:id", midilware.UserAuth, controls.DeleteCart)

		//payments route
		User.POST("/payment/cod", midilware.UserAuth, controls.CashOnDelivery)
		User.GET("/payment/razorpay", midilware.UserAuth, controls.Razorpay)
		User.GET("/payment/success", midilware.UserAuth, controls.RazorpaySuccess)
		User.GET("/success", midilware.UserAuth, controls.Success)

	}
}
