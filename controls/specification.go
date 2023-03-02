package controls

import (
	"bytes"
	"fmt"

	"text/template"

	"os/exec"

	"path/filepath"
	"strconv"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/athunlal/config"
	"github.com/athunlal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jung-kurt/gofpdf"
	"github.com/tealeg/xlsx"
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
	DB := config.DB
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
	db := config.DB
	result := db.First(&brandData)
	if result.Error != nil {
		c.JSON(500, gin.H{
			"Message": "Brand is empty",
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
	DB := config.DB

	result := DB.Model(&editbrands).Updates(models.Brand{
		BrandName: editbrands.BrandName,
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

	//binding the data from the input
	if c.Bind(&bindData) != nil {
		c.JSON(400, gin.H{
			"Bad Request": "Could not bind the JSON data",
		})
		return
	}

	//fetching the user id from the tocken
	id, err := strconv.Atoi(c.GetString("userid"))
	if err != nil {
		c.JSON(400, gin.H{
			"Error": "Error in string conversion",
		})
		return
	}

	db := config.DB

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

	//checking the produt_id and user_id  in the carts table
	err = db.Table("carts").Where("product_id = ? AND userid = ? ", bindData.Product_id, id).Select("quantity", "total_price").Row().Scan(&sum, &Price)
	fmt.Println("this is the erro : ", err)
	if err != nil {
		totalprice := productData.Price * bindData.Quantity
		cartitems := models.Cart{
			ProductId:  bindData.Product_id,
			Quantity:   bindData.Quantity,
			Price:      productData.Price,
			TotalPrice: totalprice,
			Userid:     uint(id),
		}

		//Creating the table carts
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
	totalPrice := productData.Price * totalQuantity

	//updating the quatity and the total price  to the carts
	result = db.Model(&models.Cart{}).Where("product_id = ? AND userid = ? ", bindData.Product_id, id).Updates(map[string]interface{}{"quantity": totalQuantity, "total_price": totalPrice})
	if result.Error != nil {
		c.JSON(400, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"Message": "Quantity added Successfully",
	})
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
		Product_name string
		Quantity     uint
		Total_price  uint
		Image        string
		Price        string
	}

	var datas []cartdata

	db := config.DB

	result := db.Table("carts").
		Select("products.product_name, images.image, carts.quantity, carts.price, carts.total_price").
		Joins("INNER JOIN products ON products.product_id=carts.product_id").
		Joins("INNER JOIN images ON images.product_id=carts.product_id").
		Where("userid = ?", id).
		Scan(&datas)

	if result.Error != nil {
		c.JSON(404, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}
	fmt.Println("this is the slice : ", datas)
	if len(datas) > 0 {
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

	db := config.DB
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
	pidconv := c.PostForm("product_id")
	pid, _ := strconv.Atoi(pidconv)

	db := config.DB
	var product models.Product
	result := db.First(&product, pid)
	if result.Error != nil {
		c.JSON(400, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}

	extension := filepath.Ext(imagepath.Filename)
	image := uuid.New().String() + extension
	c.SaveUploadedFile(imagepath, "./public/images"+image)

	imagedata := models.Image{
		Image:     image,
		ProductId: uint(pid),
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
	db := config.DB

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

	db := config.DB

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
	db := config.DB

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
	productId, err := strconv.Atoi(c.Query("pid"))
	if err != nil {
		c.JSON(400, gin.H{
			"Error": "Error in string conversion",
		})
	}
	db := config.DB

	//To check the product is exist or not
	var product models.Product
	result := db.First(&product, productId)
	if result.Error != nil {
		c.JSON(404, gin.H{
			"Message": "Product not exist",
		})
		return
	}
	Data := models.Wishlist{
		ProductId: uint(productId),
		Userid:    uint(userId),
	}
	result = db.Create(&Data)
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
	db := config.DB
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
			"message":  "Catogery created",
			"Catogery": createData,
		})
	} else {
		c.JSON(400, gin.H{
			"message": "Catogery already exist",
		})
	}
}

//Searching the catagery using catagery
func FilteringByCatogery(c *gin.Context) {
	cid, err := strconv.Atoi(c.Query("cid"))
	if err != nil {
		c.JSON(400, gin.H{
			"Error": "Error in string conversion ",
		})
	}
	var product models.Product
	db := config.DB
	result := db.Table("products").Select("product_name, description, price").Where("catogery_id = ? ", cid).Scan(&product).Error
	if result != nil {
		c.JSON(400, gin.H{
			"Error": "Catogery does't exist",
		})
		return
	}
	c.JSON(200, gin.H{
		"Products name": product.ProductName,
		"Discription":   product.Description,
		"Prince":        product.Price,
	})

}

//Searching the product using product name and brand name. If product name does't exist then it search using the brand name
func SearchProduct(c *gin.Context) {
	type Data struct {
		SearchValue  string
	}
	type product struct{
		Product_name string
		Description  string
		Price        float64
	}

	var userEnterData Data
	if c.Bind(&userEnterData) != nil {
		c.JSON(400, gin.H{
			"Error": "countl not bind the JSON data",
		})
	}
	var products []product
	db := config.DB
	var count int64
	result := db.Raw("SELECT product_name,description,price FROM products WHERE brand_id IN (SELECT id FROM brands WHERE brand_name LIKE ?)", "%"+userEnterData.SearchValue+"%").Scan(&products).Count(&count)
	if result.Error != nil {
		c.JSON(400, gin.H{
			"Error": "Product not exist",
		})
		return
	}

	if count <= 0 {
		result := db.Raw("SELECT * FROM products WHERE product_name LIKE ?", "%"+userEnterData.SearchValue+"%").Find(&products)
		if result.Error != nil {
			c.JSON(400, gin.H{
				"Error": result.Error.Error(),
			})
		}
		return
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

	db := config.DB
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
	result = db.Last(&oderData).Where("useridno = ? AND oder_itemid = ?", id, Oder_item.OrderId)
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
	result = db.First(&address, oderData.AddressId)
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
	err = db.Joins("JOIN oder_details ON products.product_id = oder_details.product_id").
		Where("oder_details.oder_item_id = ?", oderData.OderItemId).Find(&products).Error
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
			Product:     data.ProductName,
			Price:       data.Price,
			Description: data.Description,
		}
	}

	//spliting the date from time.time
	timeString := Payment.Date.Format("2006-01-02")

	//executing the template Invoice
	invoice := Invoice{
		Name:          user.FirstName,
		Date:          timeString,
		Email:         user.Email,
		OrderId:       oderData.OderItemId,
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

	cmd := exec.Command("wkhtmltopdf", "-", "./public/invoice.pdf")
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

//<<<<<<<<<<<<< Sales Report >>>>>>>>>>>>>>>>>>>>>>>>
func SalesReport(c *gin.Context) {

	//fetching the dates from the URL
	startDate := c.Query("startDate")
	endDateStr := c.Query("endDate")

	//converting the dates string to time.time
	fromTime, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Invalid start Date",
		})
		return
	}
	toTime, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Invalid end Date",
		})
		return
	}

	//fetching the data from the table Order details where start date to end date
	var orderDetail []models.OderDetails
	// var reportData []Report
	db := config.DB

	result := db.Preload("Product").Preload("Payment").
		Where("created_at BETWEEN ? AND ?", fromTime, toTime).
		Find(&orderDetail)

	if result.Error != nil {
		c.JSON(400, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}

	f := excelize.NewFile()

	// Create a new sheet.
	SheetName := "Sheet1"
	index := f.NewSheet(SheetName)

	// Set the value of headers
	f.SetCellValue(SheetName, "A1", "Order Date")
	f.SetCellValue(SheetName, "B1", "Order ID")
	f.SetCellValue(SheetName, "C1", "Product name")
	f.SetCellValue(SheetName, "D1", "Price")
	f.SetCellValue(SheetName, "E1", "Total Amount")
	f.SetCellValue(SheetName, "F1", "Payment method")
	f.SetCellValue(SheetName, "G1", "Payment Status")
	// Set the value of cell
	for i, report := range orderDetail {
		row := i + 2
		f.SetCellValue(SheetName, fmt.Sprintf("A%d", row), report.CreatedAt.Format("01/02/2006"))
		f.SetCellValue(SheetName, fmt.Sprintf("B%d", row), report.Oderid)
		f.SetCellValue(SheetName, fmt.Sprintf("C%d", row), report.Product.ProductName)
		f.SetCellValue(SheetName, fmt.Sprintf("D%d", row), report.Product.Price)
		f.SetCellValue(SheetName, fmt.Sprintf("E%d", row), report.Payment.Totalamount)
		f.SetCellValue(SheetName, fmt.Sprintf("F%d", row), report.Payment.PaymentMethod)
		f.SetCellValue(SheetName, fmt.Sprintf("G%d", row), report.Payment.Status)

	}

	// Set active sheet of the workbook.
	f.SetActiveSheet(index)

	// Save the Excel file with the name "test.xlsx".
	if err := f.SaveAs("./public/SalesReport.xlsx"); err != nil {
		fmt.Println(err)
	}
	CovertingExelToPdf(c)
	c.HTML(200, "SalseReport.html", gin.H{})

}

func CovertingExelToPdf(c *gin.Context) {
	// Open the Excel file
	xlFile, err := xlsx.OpenFile("./public/SalesReport.xlsx")
	if err != nil {
		fmt.Println(err)
		return
	}

	// Create a new PDF document
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "", 14)

	// Convertig each cell in the Excel file to a PDF cell
	for _, sheet := range xlFile.Sheets {
		for _, row := range sheet.Rows {
			for _, cell := range row.Cells {
				//if there is any empty cell values skiping that
				if cell.Value == "" {
					continue
				}

				pdf.Cell(40, 10, cell.Value)
			}
			pdf.Ln(-1)
		}
	}

	// Save the PDF document
	err = pdf.OutputFileAndClose("./public/SalesReport.pdf")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("PDF saved successfully.")
}

func DownloadExel(c *gin.Context) {
	c.Header("Content-Disposition", "attachment; filename=SalesReport.xlsx")
	c.Header("Content-Type", "application/xlsx")
	c.File("./public/SalesReport.xlsx")
}

func Downloadpdf(c *gin.Context) {
	c.Header("Content-Disposition", "attachment; filename=SalesReport.pdf")
	c.Header("Content-Type", "application/pdf")
	c.File("./public/SalesReport.pdf")
}

//Wallet history
func WalletHistory(c *gin.Context) {
	userid, err := strconv.Atoi(c.GetString("userid"))
	if err != nil {
		c.JSON(400, gin.H{
			"Error": err.Error(),
		})
	}
	var WalletHistory []models.WalletHistory
	db := config.DB
	result := db.Find(&WalletHistory).Where("user_id", userid)
	if result.Error != nil {
		c.JSON(400, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}
	var response []map[string]interface{}
	for _, history := range WalletHistory {
		row := map[string]interface{}{
			"Amount":          history.Amount,
			"TransactionType": history.TransctionType,
			"Date":            history.Date,
		}
		response = append(response, row)
	}

	c.JSON(200, gin.H{"data": response})
}
