package controls

import (
	"github.com/athunlal/config"
	"github.com/athunlal/models"
	"github.com/gin-gonic/gin"
)

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

// func AddBrand(c *gin.Context) {

// 	var brand models.Brand

// 	if c.Bind(&brand) != nil {
// 		c.JSON(400, gin.H{
// 			"error": "Data binding error",
// 		})
// 		return
// 	}

// 	db := config.DBconnect()
// 	result := db.Create(&brand)
// 	if result.Error != nil {
// 		c.JSON(404, gin.H{
// 			"Error": result.Error.Error(),
// 		})
// 		return
// 	}
// 	c.JSON(200, gin.H{
// 		"Message":      "Successfully Added the brand",
// 		"Product data": brand,
// 	})
// }
