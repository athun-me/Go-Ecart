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

//Admin adding the product brand
func AddBrands(c *gin.Context) {
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

//view brand
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

//Edit brand by admin
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

//Adding the product to the cart
func AddToCart(c *gin.Context) {
	type data struct {
		Product_id uint
		Quantity   uint
	}
	var bindData data
	var productData models.Product

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
		return
	}
	db := config.DBconnect()

	//checking the product is exist or not
	result := db.First(&productData, bindData.Product_id)
	if result.Error != nil {
		c.JSON(400, gin.H{
			"Message": "Product not exist",
		})
		return
	}

	//checking stock quantity
	if bindData.Quantity > productData.Stock {
		c.JSON(404, gin.H{
			"Message": "Out of Stock",
		})
		return
	}

	var sum uint
	var Price uint

	//checking the produt_id and user_id is in the carts table
	err = db.Table("carts").Where("product_id = ? AND userid = ? ", bindData.Product_id, id).Select("quantity", "totalprice").Row().Scan(&sum, &Price)
	if err != nil {
		totalprice := productData.Price * bindData.Quantity
		cartitems := models.Cart{
			Product_id: bindData.Product_id,
			Quantity:   bindData.Quantity,
			Price:      productData.Price,
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
		return
	}

	//calculatin the tottal quantity and total Price
	totalQuantity := sum + bindData.Quantity
	totalPrice := Price * totalQuantity

	//updating the quatity and the total price  to the carts
	result = db.Model(&models.Cart{}).Where("product_id = ?", bindData.Product_id).Updates(map[string]interface{}{"quantity": totalQuantity, "totalprice": totalPrice})
	if result.Error != nil {
		c.JSON(400, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"Message": "Quantity added Successfully",
	})
	return

}

//View cart items using user id
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

//Delete cart of a perticular user id
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

//Add image of the product by admin
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

//Coupon adding by admin
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

//Checking the coupon is valide or exist in the data base
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

//Coupon applying
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

//To add the wish list
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

//Adding the catogery by admin
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

//Searching the catagery using catagery
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

//Searching the product using product name and brand name. If product name does't exist then it search using the brand name
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

// creating pdf file containing invoice for show to the user
type Invoice struct {
	Name          string
	Email         string
	PaymentMethod string
	Totalamount   int64
	Date          string
	OrderId       uint
	Address       []Address
	Items         []Item
}
type Address struct {
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

type Item struct {
	Product     string
	Description string
	Qty         uint
	Price       uint
}

//Templates for creating the pdf
const invoiceTemplate = `
Order ID : {{.OrderId}}<br>
Order Date:{{.Date}} <br><hr>
Name : {{.Name}} <br>
Email: {{.Email}}<br>
<hr>
Billing Address :
{{range .Address}}

Phone number : {{.Phoneno}} <br>
House number : {{.Houseno}} <br>
Area : {{.Area}} <br>
Landmark : {{.Landmark}} <br>
City : {{.City}} <br>
Pincode : {{.Pincode}} <br>
District : {{.District}} <br>
State : {{.State}} <br>
Country : {{.Country}} <br>
{{end}}
<hr>
Payment method : {{.PaymentMethod}}<br>
<hr>
{{range .Items}}

Product :{{.Product}}  <br>
Description: {{.Description}}<br>
Price : {{.Price}}<br><br>

{{end}}
<hr><br>
Total Amount : {{.Totalamount}}<br>
`

func InvoiceF(c *gin.Context) {
	fmt.Println()

	id, err := strconv.Atoi(c.GetString("userid"))
	if err != nil {
		c.JSON(400, gin.H{
			"Error": "Error in string conversion",
		})
		return
	}

	db := config.DBconnect()
	var user models.User
	var Payment models.Payment
	var oderData models.OderDetails
	var address models.Address
	var Oder_item models.Oder_item

	//fetching the data from table Oder_item using usder id
	result := db.Last(&Oder_item).Where("useridno = ?", id)
	if result.Error != nil {
		c.JSON(400, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}

	//fetching the data from table Oder_details using userid and oder_idtemid, for fetching the oder_itemid
	result = db.Last(&oderData).Where("useridno = ? AND oder_itemid = ?", id, Oder_item.Order_id)
	if result.Error != nil {
		c.JSON(400, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}

	//Fetching the data from table users using userid
	result = db.First(&user, id)
	if result.Error != nil {
		c.JSON(400, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}

	//fetching the user address using address id from table Oder_Details
	result = db.First(&address, oderData.Address_id)
	if result.Error != nil {
		c.JSON(400, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}

	//fetching the payment detail form table Payments using userid 
	result = db.Last(&Payment, "user_id = ?", id)
	if result.Error != nil {
		c.JSON(400, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}

	//fetching the product data from table products using Oder_itemid from table Oder_item. 
	var products []models.Product
	err = db.Joins("JOIN oder_details ON products.productid = oder_details.product_id").
		Where("oder_details.oder_itemid = ?", oderData.Oder_itemid).Find(&products).Error
	if err != nil {
		c.JSON(400, gin.H{
			"Error": "somthing went wrong",
		})
		return
	}

	//To list the product details from products, product data assign to slice items 
	items := make([]Item, len(products))
	for i, data := range products {
		items[i] = Item{
			Product:     data.Productname,
			Price:       data.Price,
			Description: data.Description,
		}
	}

	//spliting the date from time.time 
	timeString := Payment.Date.Format("2006-01-02")

	//executing the template Invoice
	invoice := Invoice{
		Name:          user.Firstname,
		Date:          timeString,
		Email:         user.Email,
		OrderId:       oderData.Oder_itemid,
		PaymentMethod: Payment.PaymentMethod,
		Totalamount:   int64(Payment.Totalamount),
		Address: []Address{
			{
				Phoneno:  address.Phoneno,
				Houseno:  address.Houseno,
				Area:     address.Area,
				Landmark: address.Landmark,
				City:     address.City,
				Pincode:  address.Pincode,
				District: address.District,
				State:    address.State,
				Country:  address.Country,
			},
		},
		Items: items,
	}

	tmpl, err := template.New("invoice").Parse(invoiceTemplate)
	if err != nil {
		c.JSON(400, gin.H{
			"Error": err.Error(),
		})
		return
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, invoice)
	if err != nil {
		c.JSON(400, gin.H{
			"Error": err.Error(),
		})
		return
	}

	cmd := exec.Command("wkhtmltopdf", "-", "invoice.pdf")
	cmd.Stdin = &buf
	err = cmd.Run()
	if err != nil {
		c.JSON(400, gin.H{
			"Error": err.Error(),
		})
		return
	}

	c.HTML(200, "invoice.html", gin.H{})
}

//To download the pdf file 
func Download(c *gin.Context) {
	c.Header("Content-Disposition", "attachment; filename=invoice.pdf")
	c.Header("Content-Type", "application/pdf")
	c.File("invoice.pdf")
}
