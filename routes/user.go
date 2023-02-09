package routes

import (
	"github.com/athunlal/controls"
	midilware "github.com/athunlal/midileware"
	"github.com/gin-gonic/gin"
)

func UserRouts(c *gin.Engine) {
	User := c.Group("/user")
	{

		//User rountes
		User.POST("/login", controls.UesrLogin)
		User.POST("/signup", controls.UserSignUP)
		User.POST("/signup/otpvalidate", controls.OtpValidation)
		User.GET("/logout", midilware.UserAuth, controls.UserSignout)
		User.GET("/validate", midilware.UserAuth, controls.Validate)

		//User profile routes
		User.GET("/viewprofile", midilware.UserAuth, controls.ShowUserDetails)
		User.POST("/addaddress", midilware.UserAuth, controls.Addaddress)
		User.PUT("/editaddress/:id", midilware.UserAuth, controls.EditUserAddress)
		User.GET("/searchaddress/:id", midilware.UserAuth, controls.ShowAddress)
		User.POST("/userchangepassword", midilware.UserAuth, controls.UserChangePassword)
		User.PUT("/updatepassword", midilware.UserAuth, controls.Updatepassword)
		User.PUT("/editprofile", midilware.UserAuth, controls.EditUserProfilebyUser)
		User.GET("/wishlist/:id", midilware.UserAuth, controls.Wishlist)

		//User product managment
		User.GET("/viewbrand", midilware.UserAuth, controls.ViewBrand)
		User.GET("/search", midilware.UserAuth, controls.SearchProduct)
		User.GET("/viewproducts", midilware.UserAuth, controls.ViewProducts)

		//User carts routes
		User.GET("/viewcart", midilware.UserAuth, controls.ViewCart)
		User.POST("/addtocart", midilware.UserAuth, controls.AddToCart)
		User.GET("/fileterbycatogery/:id", midilware.UserAuth, controls.FilteringByCatogery)
		User.GET("/checkout", midilware.UserAuth, controls.CheckOut)
		User.DELETE("/deletecart/:id", midilware.UserAuth, controls.DeleteCart)

		//Oder managements by user
		User.GET("/showoder", midilware.UserAuth, controls.ShowOder)
		User.GET("/return", midilware.UserAuth, controls.ReturnOder)
		User.GET("/canceloder", midilware.UserAuth, controls.CancelOder)

		//Coupon management
		User.POST("/applycoupon", midilware.UserAuth, controls.Applycoupon)
		User.POST("/checkcoupon", midilware.UserAuth, controls.CheckCoupon)

		//Forgot Password
		User.PUT("/forgotpasswordotpvalidation", midilware.UserAuth, controls.ForgotPasswordOtpValidation)
		User.PUT("/forgotpassword", midilware.UserAuth, controls.ForgotPassword)
		User.POST("/changepassword", midilware.UserAuth, controls.ChangePassword)

		//payments route
		User.GET("/payment/cashOnDelivery", midilware.UserAuth, controls.CashOnDelivery)
		User.GET("/payment/razorpay", midilware.UserAuth, controls.Razorpay)
		User.GET("/payment/success", midilware.UserAuth, controls.RazorpaySuccess)
		User.GET("/success", midilware.UserAuth, controls.Success)

	}
}
