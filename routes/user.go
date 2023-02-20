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

		//User profile routes
		User.GET("/profile/viewprofile", middlereware.UserAuth, controls.ShowUserDetails)
		User.POST("/prfile/userchangepassword", middlereware.UserAuth, controls.UserChangePassword)
		User.PUT("/prfile/userchangepassword/updatepassword", middlereware.UserAuth, controls.Updatepassword)
		User.PUT("/profile/editprofile", middlereware.UserAuth, controls.EditUserProfilebyUser)
		User.POST("/addaddress", middlereware.UserAuth, controls.Addaddress)
		User.PUT("/address/editaddress/:id", middlereware.UserAuth, controls.EditUserAddress)
		User.GET("/address/searchaddress/:id", middlereware.UserAuth, controls.ShowAddress)

		//User product managment
		User.GET("/product/wishlist/:id", middlereware.UserAuth, controls.Wishlist)
		User.GET("/product/brand/viewbrand", middlereware.UserAuth, controls.ViewBrand)
		User.GET("/product/search", middlereware.UserAuth, controls.SearchProduct)
		User.GET("/product/viewproducts", middlereware.UserAuth, controls.ViewProducts)

		//User carts routes
		User.GET("/profile/viewcart", middlereware.UserAuth, controls.ViewCart)
		User.POST("/profile/addtocart", middlereware.UserAuth, controls.AddToCart)
		User.GET("/profile/fileterbycatogery/:id", middlereware.UserAuth, controls.FilteringByCatogery)
		User.DELETE("/deletecart/:id", middlereware.UserAuth, controls.DeleteCart)

		//Oder managements by user
		User.GET("/order/showoder", middlereware.UserAuth, controls.ShowOder)
		User.GET("/order/showoder/return", middlereware.UserAuth, controls.ReturnOrderByUser)
		User.GET("/order/showoder/canceloder", middlereware.UserAuth, controls.CancelOrder)
		User.GET("/cart/checkout", middlereware.UserAuth, controls.CheckOut)

		//Coupon management
		User.POST("/cart/checkout/applycoupon", middlereware.UserAuth, controls.Applycoupon)
		User.POST("/checkcoupon", middlereware.UserAuth, controls.CheckCoupon)

		//Forgot Password
		User.PUT("/forgotpassword", middlereware.UserAuth, controls.GenerateOtpForForgotPassword)
		User.POST("/forgotpassword/changepassword", middlereware.UserAuth, controls.ChangePassword)

		//payments route
		User.GET("/payment/cashOnDelivery", middlereware.UserAuth, controls.CashOnDelivery)
		User.GET("/payment/razorpay", middlereware.UserAuth, controls.Razorpay)
		User.GET("/payment/success", middlereware.UserAuth, controls.RazorpaySuccess)
		User.GET("/success", middlereware.UserAuth, controls.Success)

		//Invoice download
		User.GET("/invoice", middlereware.UserAuth, controls.InvoiceF)
		User.GET("/invoice/download", middlereware.UserAuth, controls.Download)

	}
}
