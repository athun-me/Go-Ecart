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
