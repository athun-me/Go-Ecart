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
	err := c.Bind(&product)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Data binding error",
		})
		fmt.Println(err)
		return
	}

	db := config.DB
	var count int64
	result := db.Find(&product, "product_name = ?", product.ProductName).Count(&count)
	if result.Error != nil {
		c.JSON(404, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}
	if count == 0 {
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
	} else {
		c.JSON(400, gin.H{
			"Message": "Product already exist",
		})
	}
}

//>>>>>>>>>>>>>>>>> View products <<<<<<<<<<<<<<<<<<<<<
func ViewProducts(c *gin.Context) {
	limit, _ := strconv.Atoi(c.Query("limit"))
	offset, _ := strconv.Atoi(c.Query("offset"))
	type datas struct {
		Product_name string
		Description string
		Stock       string
		Price       string
		Brand_name   string
	}
	var products datas

	db := config.DB
	query := "SELECT products.product_name, products.description, products.stock, products.price, brands.brand_name FROM products LEFT JOIN brands ON products.brand_id=brands.id  GROUP BY products.product_id, brands.brand_name"

	if limit != 0 || offset != 0 {
		if limit == 0 {
			query = fmt.Sprintf("%s OFFSET %d", query, offset)
		} else if offset == 0 {
			query = fmt.Sprintf("%s LIMIT %d", query, limit)
		} else {
			query = fmt.Sprintf("%s LIMIT %d OFFSET %d", query, limit, offset)
		}
	}
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
