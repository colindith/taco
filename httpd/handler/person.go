package handler

import (
	"fmt"
	// "net/http"
	// "taco/packages/gredis"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type Person struct {
	ID        uint   `json:"id"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	City      string `json:"city"`
}

func DeletePerson(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Params.ByName("id")
		var person Person
		d := db.Where("id = ?", id).Delete(&person)
		fmt.Println(d)
		c.JSON(200, gin.H{"id #" + id: "deleted"})
	}
}
func UpdatePerson(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var person Person
		id := c.Params.ByName("id")
		if err := db.Where("id = ?", id).First(&person).Error; err != nil {
			c.AbortWithStatus(404)
			fmt.Println(err)
		}
		c.BindJSON(&person)
		db.Save(&person)
		c.JSON(200, person)
	}
}
func CreatePerson(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var person Person
		// var cache string
		c.BindJSON(&person)
		db.Create(&person)
		c.JSON(200, person)
	}
}
func GetPerson(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Params.ByName("id")
		var person Person
		if err := db.Where("id = ?", id).First(&person).Error; err != nil {
			c.AbortWithStatus(404)
			fmt.Println(err)
		} else {
			c.JSON(200, person)
		}
	}
}
func GetPeople(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var people []Person
		if err := db.Find(&people).Error; err != nil {
			c.AbortWithStatus(404)
			fmt.Println(err)
		} else {
			c.JSON(200, people)
		}
	}
}
