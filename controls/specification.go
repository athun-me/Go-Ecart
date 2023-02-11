package controls

import (
	"bytes"
	"fmt"
	"text/template"

	"os/exec"

	"path/filepath"
	"strconv"
	"time"

	"github.com/athunlal/config"
	"github.com/athunlal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

//>>>>>> Add brand <<<<<<<<<<<<<<<<<<<<
func AddBrands(c *gin.Context) {
	fmt.Println("><")
	var addbrand models.Brand
	if c.Bind(&addbrand) != nil {
		c.JSON(400, gin.H{
			"Error": "Could not bind JSON data",
		})
		return
	}
	DB := config.DBconnect()
	result := DB.Create(&addbrand)
	if result.Error != nil {
		c.JSON(500, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"Message":       "New Brand added Successfully",
		"Brand details": addbrand,
	})
}

// >>>>>>>>>> view brand <<<<<<<<<<<<<<<<<<<<<<<<<<<<<
func ViewBrand(c *gin.Context) {
	var brandData []models.Brand
	db := config.DBconnect()
	result := db.Find(&brandData)
	if result.Error != nil {
		c.JSON(500, gin.H{
			"Status":  "False",
			"Message": "Could not retrieve the brands",
		})
		return
	}
	c.JSON(200, gin.H{
		"Brands data": brandData,
	})
}

//>>>>>>>>>>>>>> Edit brand <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<
func EditBrand(c *gin.Context) {
	bid := c.Param("id")
	id, err := strconv.Atoi(bid)
	if err != nil {
		c.JSON(400, gin.H{
			"Error": "Error in string conversion",
		})
	}
	var editbrands models.Brand
	if c.Bind(&editbrands) != nil {
		c.JSON(400, gin.H{
			"Error": "Error in binding the JSON data",
		})
		return
	}
	editbrands.ID = uint(id)
	DB := config.DBconnect()

	result := DB.Model(&editbrands).Updates(models.Brand{
		Brandname: editbrands.Brandname,
	})

	if result.Error != nil {
		c.JSON(404, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"Message": "Successfully updated the Brand",
	})
}

//>>>>>>>>>>>> Add to cart <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<
func AddToCart(c *gin.Context) {
	type data struct {
		Product_id uint
		Quantity   uint
	}
	var bindData data
	var cartdata models.Cart
	var productdata models.Product
	if c.Bind(&bindData) != nil {
		c.JSON(400, gin.H{
			"Bad Request": "Could not bind the JSON data",
		})
		return
	}

	id, err := strconv.Atoi(c.GetString("userid"))
	if err != nil {
		c.JSON(400, gin.H{
			"Error": "Error in string conversion",
		})
	}

	db := config.DBconnect()
	var count int64

	//fetching the table products for checking stocks
	db.Table("products").Select("stock, price").Where("productid = ?", bindData.Product_id).Scan(&productdata)
	if bindData.Quantity >= productdata.Stock {
		c.JSON(404, gin.H{
			"Message": "Out of Stock",
		})
		return
	}

	//fetching the table carts for checking the product_id is exist
	db.Find(&cartdata, "userid = ?", id).Count(&count)

	if count > 0 {
		var sum uint

		//fetching the quantity form carts
		db.Table("carts").Where("product_id = ? AND userid = ? ", bindData.Product_id, id).Select("SUM(quantity)").Row().Scan(&sum)
		totalQuantity := sum + bindData.Quantity

		//updating the quatity to the carts
		db.Model(&cartdata).Where("product_id = ?", bindData.Product_id).Update("quantity", totalQuantity)
		c.JSON(200, gin.H{
			"Message": "Quantity added Successfully",
		})
		return
	}

	totalprice := productdata.Price * bindData.Quantity
	cartitems := models.Cart{
		Product_id: bindData.Product_id,
		Quantity:   bindData.Quantity,
		Price:      productdata.Price,
		Totalprice: totalprice,
		Userid:     uint(id),
	}
	result := db.Create(&cartitems)
	if result.Error != nil {
		c.JSON(400, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"Message": "Added to the Cart Successfull",
	})
}

//>>>>>>>>>>>>>>> View Products <<<<<<<<<<<<<<<<<<<<<<<

func ViewCart(c *gin.Context) {
	id, err := strconv.Atoi(c.GetString("userid"))
	if err != nil {
		c.JSON(400, gin.H{
			"Error": "Error in string conversion",
		})
	}
	type cartdata struct {
		Productname string
		Quantity    uint
		Totalprice  uint
		Image       string
		Price       string
	}
	var datas []cartdata
	db := config.DBconnect()
	result := db.Table("carts").Select("products.productname, carts.quantity, carts.price, carts.totalprice").Joins("INNER JOIN products ON products.productid=carts.product_id").Where("userid = ?", id).Scan(&datas)
	if result.Error != nil {
		c.JSON(404, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}
	if datas != nil {
		c.JSON(200, gin.H{
			"Cart Items": datas,
		})
	} else {
		c.JSON(404, gin.H{
			"Message": "Cart is empty",
		})
	}
}

//>>>>>>>>>>>>>Remove cart <<<<<<<<<<<<<<<<<<<<<
func DeleteCart(c *gin.Context) {
	id := c.Param("id")
	userid, err := strconv.Atoi(c.GetString("userid"))
	if err != nil {
		c.JSON(400, gin.H{
			"Error": "Error in string conversion",
		})
	}

	db := config.DBconnect()
	result := db.Exec("delete from carts where id= ? AND userid = ?", id, userid)
	count := result.RowsAffected
	if count == 0 {
		c.JSON(400, gin.H{
			"Message": "Cart not exist",
		})
		return
	}
	if result.Error != nil {
		c.JSON(400, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"Cart Items": "Delete successfully",
	})
}

//>>>>>>>>> Add Image <<<<<<<<<<<<<<<<<<<<<<<

func AddImages(c *gin.Context) {
	imagepath, _ := c.FormFile("image")
	extension := filepath.Ext(imagepath.Filename)
	image := uuid.New().String() + extension
	c.SaveUploadedFile(imagepath, "./public/images"+image)
	pidconv := c.PostForm("product_id")
	pid, _ := strconv.Atoi(pidconv)
	db := config.DBconnect()

	var product models.Product
	result := db.First(&product, pid)
	if result.Error != nil {
		c.JSON(400, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}

	imagedata := models.Image{
		Image:      image,
		Product_id: uint(pid),
	}
	result = db.Create(&imagedata)
	if result.Error != nil {
		c.JSON(400, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"Message": "Image Added Successfully",
	})
}

// >>>>>>>>>>>>>>>> coupon <<<<<<<<<<<<<<<<<<<
func AddCoupon(c *gin.Context) {

	type data struct {
		CouponCode    string
		Year          uint
		Month         uint
		Day           uint
		DiscountPrice float64
		Expired       time.Time
	}

	var userEnterData data
	var couponData []models.Coupon
	db := config.DBconnect()

	if c.Bind(&userEnterData) != nil {
		c.JSON(400, gin.H{
			"Error": "Could not bind the JSON data",
		})
		return
	}
	specificTime := time.Date(int(userEnterData.Year), time.Month(userEnterData.Month), int(userEnterData.Day), 0, 0, 0, 0, time.UTC)

	userEnterData.Expired = specificTime
	var count int64
	result := db.First(&couponData, "coupon_code = ?", userEnterData.CouponCode).Count(&count)
	if result.Error != nil {
		c.JSON(404, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}
	if count == 0 {
		Data := models.Coupon{
			CouponCode:    userEnterData.CouponCode,
			DiscountPrice: userEnterData.DiscountPrice,
			Expired:       specificTime,
		}
		result := db.Create(&Data)
		if result.Error != nil {
			c.JSON(400, gin.H{
				"Error": result.Error.Error(),
			})
		}
		c.JSON(200, gin.H{
			"message": userEnterData,
		})
	} else {
		c.JSON(400, gin.H{
			"message": "Coupon already exist",
		})
	}

}

//>>>>>>>>>>>>>> Check coupon <<<<<<<<<<<<<<<<<<<<
func CheckCoupon(c *gin.Context) {
	type data struct {
		Coupon string
	}
	var coupon models.Coupon
	var userEnterData data

	if c.Bind(&userEnterData) != nil {
		c.JSON(400, gin.H{
			"Error": "Could not bind the JSON data",
		})
		return
	}

	db := config.DBconnect()

	var count int64
	result := db.Find(&coupon, "coupon_code = ?", userEnterData.Coupon).Count(&count)
	if result.Error != nil {
		c.JSON(400, gin.H{
			"Error": result.Error.Error(),
		})
	}

	if count == 0 {
		c.JSON(400, gin.H{
			"message": "Coupon not exist",
		})
		return
	}
	currentTime := time.Now()
	expiredData := coupon.Expired

	if currentTime.Before(expiredData) {
		c.JSON(200, gin.H{
			"message": "Coupon valide",
		})
	} else if currentTime.After(expiredData) {
		c.JSON(400, gin.H{
			"message": "Coupon expired",
		})
	}
}

//>>>>>>>>>>>>> Applying coupon <<<<<<<<<<<<<<<<<<<<<<<<<

func Applycoupon(c *gin.Context) {
	id, err := strconv.Atoi(c.GetString("userid"))
	if err != nil {
		c.JSON(400, gin.H{
			"Error": "Error while string conversion",
		})
	}
	type data struct {
		Coupon   string
		Name     string
		Phoneno  string
		Houseno  string
		Area     string
		Landmark string
		City     string
		Pincode  string
		District string
		State    string
		Country  string
	}
	var userEnterData data
	var coupon models.Coupon
	var discountPercentage float64
	var discountPrice float64
	if c.Bind(&userEnterData) != nil {
		c.JSON(400, gin.H{
			"Error": "Could not bind the JSON data",
		})
		return
	}
	db := config.DBconnect()

	//checking coupon is existig or not
	var count int64
	result := db.Find(&coupon, "coupon_code = ?", userEnterData.Coupon).Count(&count)
	if result.Error != nil {
		c.JSON(400, gin.H{
			"Error": result.Error.Error(),
		})
	}
	if count == 0 {
		c.JSON(400, gin.H{
			"message": "Coupon not exist",
		})
	} else {
		currentTime := time.Now()
		expiredData := coupon.Expired

		if currentTime.Before(expiredData) {
			c.JSON(200, gin.H{
				"message": "Coupon valide",
			})
			discountPercentage = coupon.DiscountPrice
		} else if currentTime.After(expiredData) {
			c.JSON(400, gin.H{
				"message": "Coupon expired",
			})
		}
	}

	//fetching the cart details from the table carts
	ViewCart(c)

	//fetching the data from table addresses
	result = db.Raw("SELECT name, phoneno, houseno, area, landmark, city, pincode,district, state, country FROM addresses WHERE defaultadd = true AND userid = ?", id).Scan(&userEnterData)
	if result.Error != nil {
		c.JSON(404, gin.H{
			"Error": result,
		})
		return
	}
	c.JSON(200, gin.H{
		"Default Address of user": userEnterData,
	})

	//fetching and calculatin the total amount of the cart products
	var totalPrice float64
	result1 := db.Table("carts").Where("userid = ?", id).Select("SUM(totalprice)").Scan(&totalPrice).Error

	//calculating the discount amount
	discountPrice = (discountPercentage / 100) * totalPrice
	totalPriceAfterDeduct := totalPrice - discountPrice

	if result1 != nil {
		c.JSON(400, gin.H{
			"Error": "Can not fetch total amount",
		})
		return
	}
	c.JSON(200, gin.H{
		"Price":        totalPrice,
		"Discount":     discountPrice,
		"Total amount": totalPriceAfterDeduct,
	})
}

//>>>>>>>>>>>>>> wish list  <<<<<<<<<<<<<<<<<<<
func Wishlist(c *gin.Context) {
	userId, err := strconv.Atoi(c.GetString("userid"))
	if err != nil {
		c.JSON(400, gin.H{
			"Error": "Error in string conversion",
		})
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{
			"Error": "Error in string conversion",
		})
	}

	db := config.DBconnect()
	Data := models.Wishlist{
		Product_id: uint(id),
		Userid:     uint(userId),
	}
	result := db.Create(&Data)
	if result.Error != nil {
		c.JSON(400, gin.H{
			"Error": result.Error.Error(),
		})
	}
	c.JSON(200, gin.H{
		"message": "Wish list added sucessfully",
	})
}

//>>>>>>>>>>>>>>>>>> Add catogeries <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<

func AddCatogeries(c *gin.Context) {

	type Data struct {
		CategoryName string
	}
	var category Data
	var CatagoryData models.Catogery
	if c.Bind(&category) != nil {
		c.JSON(400, gin.H{
			"Error": "countl not bind the JSON data",
		})
	}
	db := config.DBconnect()
	var count int64
	result := db.Find(&CatagoryData, "catogery_name = ?", category.CategoryName).Count(&count)
	if result.Error != nil {
		c.JSON(400, gin.H{
			"Error": result.Error.Error(),
		})
	}
	if count == 0 {
		createData := models.Catogery{
			CatogeryName: category.CategoryName,
		}
		result = db.Create(&createData)
		if result.Error != nil {
			c.JSON(400, gin.H{
				"Error": result.Error.Error(),
			})
		}
		c.JSON(200, gin.H{
			"message": "Catogery created",
		})
	} else {
		c.JSON(400, gin.H{
			"message": "Catogery already exist",
		})
	}
}

//>>>>>>>>>>Search by catogery <<<<<<<<<<<<<<<<<<<<<<<<<<<

func FilteringByCatogery(c *gin.Context) {
	id := c.Param("id")

	var products []models.Product
	db := config.DBconnect()
	result := db.Where("catogery_id = ?", id).Find(&products)
	if result.Error != nil {
		c.JSON(400, gin.H{
			"Error": result.Error.Error(),
		})
	}
	c.JSON(200, gin.H{
		"products": products,
	})

}

//>>>>>>>>>>>>>>>>>> Search <<<<<<<<<<<<<<<<<<<<<

func SearchProduct(c *gin.Context) {
	type Data struct {
		SearchValue string
	}
	var userEnterData Data
	if c.Bind(&userEnterData) != nil {
		c.JSON(400, gin.H{
			"Error": "countl not bind the JSON data",
		})
	}
	var products []models.Product
	db := config.DBconnect()
	var count int64
	result := db.Raw("SELECT * FROM products WHERE brand_id (SELECT id FROM brands WHERE brandname LIKE ?)", "%"+userEnterData.SearchValue+"%").Scan(&products).Count(&count)
	if result.Error != nil {
		c.JSON(400, gin.H{
			"Error": result.Error.Error(),
		})
	}

	if count <= 0 {
		result := db.Raw("SELECT * FROM products WHERE productname LIKE ?", "%"+userEnterData.SearchValue+"%").Find(&products)
		if result.Error != nil {
			c.JSON(400, gin.H{
				"Error": result.Error.Error(),
			})
		}
	}

	if count == 0 {
		c.JSON(400, gin.H{
			"Message": "Product not exist",
		})
		return
	}
	c.JSON(200, gin.H{
		"products": products,
	})
}

//>>>>>>>>>> Oder Details <<<<<<<<<<<<<<<<

func OderDetails(c *gin.Context) {
	userId, err := strconv.Atoi(c.GetString("userid"))
	if err != nil {
		c.JSON(400, gin.H{
			"Error": "Error in string conversion",
		})
	}

	var UserAddress models.Address
	var UserPayment models.Payment
	var UserCart []models.Cart

	db := config.DBconnect()
	result := db.Find(&UserAddress, "userid = ? AND defaultadd = true", userId)
	if result.Error != nil {
		c.JSON(400, gin.H{
			"Error": result.Error.Error(),
		})
	}
	result = db.Find(&UserPayment, "user_id = ?", userId)
	if result.Error != nil {
		c.JSON(400, gin.H{
			"Error": result.Error.Error(),
		})
	}
	result = db.Find(&UserCart, "userid = ?", userId)
	if result.Error != nil {
		c.JSON(400, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}

	for _, UserCart := range UserCart {
		OderDetails := models.OderDetails{
			Userid:     uint(userId),
			Address_id: UserAddress.Addressid,
			Paymentid:  UserPayment.Payment_id,
			Product_id: UserCart.Product_id,
			Quantity:   UserCart.Quantity,
			Status:     "Pending",
		}

		result = db.Create(&OderDetails)
		if result.Error != nil {
			c.JSON(400, gin.H{
				"Error": result.Error.Error(),
			})
			return
		}
	}
	c.JSON(200, gin.H{
		"Message": "Oder Added succesfully",
	})
}

//>>>>>>>>>> Show oder <<<<<<<<<<<<<<<<<<<
func ShowOder(c *gin.Context) {
	userId, err := strconv.Atoi(c.GetString("userid"))
	if err != nil {
		c.JSON(400, gin.H{
			"Error": "Error in string conversion",
		})
		return
	}
	var userOder []models.OderDetails
	var products []models.Product
	// var userProduct models.Product

	db := config.DBconnect()
	result := db.Find(&userOder, "userid = ?", userId)
	if result.Error != nil {
		c.JSON(400, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}

	for _, order := range userOder {

		db.Find(&products, "productid = ? ", order.Product_id)

		c.JSON(200, gin.H{
			"Product name ": products[0].Productname,
			"Price":         products[0].Price,
			"Description":   products[0].Description,
			"Quantity":      userOder[0].Quantity,
		})
	}
}

//>>>>>>>>>>>>>>< Cancel Oder <<<<<<<<<<<<<<<<<<<<
func CancelOder(c *gin.Context) {
	userid, err := strconv.Atoi(c.GetString("userid"))
	if err != nil {
		c.JSON(400, gin.H{
			"Error": "Error in string conversion",
		})
	}
	var oder models.OderDetails
	db := config.DBconnect()
	result := db.Model(&oder).Where("userid = ?", userid).Update("status", "Canceled")
	if result.Error != nil {
		c.JSON(400, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"Massage": "oder canceld",
	})
}

//>>>>>>>>>>>>>>< Retrun Oder <<<<<<<<<<<<<<<<<<<<
func ReturnOder(c *gin.Context) {
	userid, err := strconv.Atoi(c.GetString("userid"))
	if err != nil {
		c.JSON(400, gin.H{
			"Error": "Error in string conversion",
		})
	}
	var oder models.OderDetails
	db := config.DBconnect()
	result := db.Model(&oder).Where("userid = ?", userid).Update("status", "Product return")
	if result.Error != nil {
		c.JSON(400, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"Massage": "Product Return",
	})
}

//>>>>>>>>> Invoice download <<<<<<<<<<<<<<<<<<
type Invoice struct {
	Name          string
	Email         string
	PaymentMethod string
	Totalamount   int64
	Items         []Item
}

type Item struct {
	Description string
	Amount      float64
}

const invoiceTemplate = `
Invoice for {{.Name}} <br>
Invoice Date: {{.Email}}

{{range .Items}}
Description: {{.Description}}
Amount: {{.Amount}}
{{end}}
`

func InvoiceF(c *gin.Context) {
	db := config.DBconnect()
	var user models.User
	db.Find(&user)

	invoice := Invoice{
		Name:  user.Firstname,
		Email: user.Email,
		// Items: []Item{
		// 	{Description: "Item 1", Amount: 100.0},
		// },
	}

	tmpl, err := template.New("invoice").Parse(invoiceTemplate)
	if err != nil {
		panic(err)
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, invoice)
	if err != nil {
		panic(err)
	}

	cmd := exec.Command("wkhtmltopdf", "-", "invoice.pdf")
	cmd.Stdin = &buf
	err = cmd.Run()
	if err != nil {
		panic(err)
	}

	defer buf.Reset()
	c.HTML(200, "invoice.html", gin.H{})
}
