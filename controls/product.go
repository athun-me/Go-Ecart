package controls

import (
	"fmt"
	"strconv"

	"github.com/athunlal/config"
	"github.com/athunlal/models"
	"github.com/gin-gonic/gin"
)

//>>>>>>>>>>>>>> Add products <<<<<<<<<<<<<<<<<<<<<<<<<<
func AddProduct(c *gin.Context) {
	var product models.Product

	if c.Bind(&product) != nil {
		c.JSON(400, gin.H{
			"error": "Data binding error",
		})
		return
	}
	db := config.DBconnect()
	result := db.Create(&product)

	if result.Error != nil {
		c.JSON(404, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"Message":      "Successfully Added the Product",
		"Product data": product,
	})
}

//>>>>>>>>>>>>>>>>> View products <<<<<<<<<<<<<<<<<<<<<
func ViewProducts(c *gin.Context) {
	limit, _ := strconv.Atoi(c.Query("limit"))
	offset, _ := strconv.Atoi(c.Query("offset"))
	type datas struct {
		Productname string
		Description string
		Stock       string
		Price       string
		Brandname   string
	}
	var products datas

	db := config.DBconnect()
	query := "SELECT products.productname, products.description, products.stock, products.price, brands.brandname FROM products LEFT JOIN brands ON products.brandid=brands.id  GROUP BY products.productid, brands.brandname"

	if limit != 0 || offset != 0 {
		if limit == 0 {
			query = fmt.Sprintf("%s OFFSET %d", query, offset)
		} else if offset == 0 {
			query = fmt.Sprintf("%s LIMIT %d", query, limit)
		} else {
			query = fmt.Sprintf("%s LIMIT %d OFFSET %d", query, limit, offset)
		}
	}
	fmt.Println(query)
	result := db.Raw(query).Scan(&products)
	if result.Error != nil {
		c.JSON(404, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"Products": products,
	})
}

//>>>>>>>>>>>> Add to cart <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<
func AddToCart(c *gin.Context) {

	var cartdata models.Cart
	var productdata models.Product
	if c.Bind(&cartdata) != nil {
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
	DB := config.DBconnect()
	DB.Table("products").Select("stock, price").Where("productid = ?", cartdata.Product_id).Scan(&productdata)

	if cartdata.Quantity >= productdata.Stock {
		c.JSON(404, gin.H{
			"Message": "Out of Stock",
		})
		return
	}
	totalprice := productdata.Price * cartdata.Quantity
	cartitems := models.Cart{
		Product_id: cartdata.Product_id,
		Quantity:   cartdata.Quantity,
		Price:      productdata.Price,
		Totalprice: totalprice,
		Cartid:     uint(id),
	}
	result := DB.Create(&cartitems)
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
	result := db.Table("carts").Select("products.productname, carts.quantity, carts.price, carts.totalprice").Joins("INNER JOIN products ON products.productid=carts.product_id").Where("cartid = ?", id).Scan(&datas)
	if result.Error != nil {
		c.JSON(404, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"Cart Items": datas,
	})
}

//>>>>>>>>>>>>>Remove cart <<<<<<<<<<<<<<<<<<<<<
func RemoveCart(c *gin.Context) {

	var cartData models.Cart
	if c.Bind(&cartData) != nil {
		c.JSON(400, gin.H{
			"Bad Request": "Could not bind the JSON data",
		})
		return
	}
	db := config.DBconnect()
	result := db.Exec("delete from carts where id= ?", cartData.ID)
	if cartData.ID == 0 {
		c.JSON(400, gin.H{
			"Bad Request": "Cart not exist",
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
