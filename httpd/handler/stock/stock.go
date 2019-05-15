package stock

import (
	"fmt"
	"time"

	// "net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type Product struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type Price struct {
	ID      uint   `json:"id"`
	product string `json:"product"`
	time    time.Time
	open    float64
	high    float64
	low     float64
	close   float64
	volume  uint
}

func DeleteProduct(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Params.ByName("id")
		var product Product
		d := db.Where("id = ?", id).Delete(&product)
		fmt.Println(d)
		c.JSON(200, gin.H{"id #" + id: "deleted"})
	}
}
func UpdateProduct(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var product Product
		id := c.Params.ByName("id")
		if err := db.Where("id = ?", id).First(&product).Error; err != nil {
			c.AbortWithStatus(404)
			fmt.Println(err)
		}
		c.BindJSON(&product)
		db.Save(&product)
		c.JSON(200, product)
	}
}
func CreateProduct(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var product Product
		c.BindJSON(&product)
		db.Create(&product)
		c.JSON(200, product)
	}
}
func GetProduct(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Params.ByName("id")
		var product Product
		if err := db.Where("id = ?", id).First(&product).Error; err != nil {
			c.AbortWithStatus(404)
			fmt.Println(err)
		} else {
			c.JSON(200, product)
		}
	}
}
func GetALlProduct(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var products []Product
		if err := db.Find(&products).Error; err != nil {
			c.AbortWithStatus(404)
			fmt.Println(err)
		} else {
			c.JSON(200, products)
		}
	}
}
