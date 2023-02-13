package routes

import (
	"github.com/athunlal/controls"
	"github.com/athunlal/middlereware"
	"github.com/gin-gonic/gin"
)

func UserRouts(c *gin.Engine) {
	User := c.Group("/user")
	{

		//User rountes
		User.POST("/login", controls.UesrLogin)
		User.POST("/signup", controls.UserSignUP)
		User.POST("/signup/otpvalidate", controls.OtpValidation)
		User.GET("/logout", middlereware.UserAuth, controls.UserSignout)
		User.GET("/validate", middlereware.UserAuth, controls.Validate)

		//User profile routes
		User.GET("/viewprofile", middlereware.UserAuth, controls.ShowUserDetails)
		User.POST("/addaddress", middlereware.UserAuth, controls.Addaddress)
		User.PUT("/editaddress/:id", middlereware.UserAuth, controls.EditUserAddress)
		User.GET("/searchaddress/:id", middlereware.UserAuth, controls.ShowAddress)
		User.POST("/userchangepassword", middlereware.UserAuth, controls.UserChangePassword)
		User.PUT("/updatepassword", middlereware.UserAuth, controls.Updatepassword)
		User.PUT("/editprofile", middlereware.UserAuth, controls.EditUserProfilebyUser)
		User.GET("/wishlist/:id", middlereware.UserAuth, controls.Wishlist)

		//User product managment
		User.GET("/viewbrand", middlereware.UserAuth, controls.ViewBrand)
		User.GET("/search", middlereware.UserAuth, controls.SearchProduct)
		User.GET("/viewproducts", middlereware.UserAuth, controls.ViewProducts)

		//User carts routes
		User.GET("/viewcart", middlereware.UserAuth, controls.ViewCart)
		User.POST("/addtocart", middlereware.UserAuth, controls.AddToCart)
		User.GET("/fileterbycatogery/:id", middlereware.UserAuth, controls.FilteringByCatogery)
		User.GET("/checkout", middlereware.UserAuth, controls.CheckOut)
		User.DELETE("/deletecart/:id", middlereware.UserAuth, controls.DeleteCart)

		//Oder managements by user
		User.GET("/showoder", middlereware.UserAuth, controls.ShowOder)
		User.GET("/return", middlereware.UserAuth, controls.ReturnOder)
		User.GET("/canceloder", middlereware.UserAuth, controls.CancelOder)

		//Coupon management
		User.POST("/applycoupon", middlereware.UserAuth, controls.Applycoupon)
		User.POST("/checkcoupon", middlereware.UserAuth, controls.CheckCoupon)

		//Forgot Password
		User.PUT("/forgotpassword", middlereware.UserAuth, controls.GenerateOtpForForgotPassword)
		User.POST("/changepassword", middlereware.UserAuth, controls.ChangePassword)

		//payments route
		User.GET("/payment/cashOnDelivery", middlereware.UserAuth, controls.CashOnDelivery)
		User.GET("/payment/razorpay", middlereware.UserAuth, controls.Razorpay)
		User.GET("/payment/success", middlereware.UserAuth, controls.RazorpaySuccess)
		User.GET("/success", middlereware.UserAuth, controls.Success)

		User.GET("/invoice", middlereware.UserAuth, controls.InvoiceF)
		User.GET("/invoice/download",middlereware.UserAuth, controls.Download)

	}
}
