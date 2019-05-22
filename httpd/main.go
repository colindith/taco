package main

import (
	"fmt"
	"io/ioutil"
	"taco/httpd/handler"
	"taco/httpd/handler/stock"
	"taco/packages/gredis"
	"taco/platform/newsfeed"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var db *gorm.DB
var err error

func init() {
	gredis.Setup()
}

func HomePage(c *gin.Context) {

	c.JSON(200, gin.H{
		"message": "Hello World ~~",
	})
}

func PostHomePage(c *gin.Context) {
	body := c.Request.Body
	value, err := ioutil.ReadAll(body)
	if err != nil {
		fmt.Println(err.Error())
	}
	c.JSON(200, gin.H{
		"message": string(value),
	})
}

func QueryStrings(c *gin.Context) {
	name := c.Query("name")
	age := c.Query("age")
	c.JSON(200, gin.H{
		"name": name,
		"age":  age,
	})
}

func PathParameters(c *gin.Context) {
	name := c.Param("name")
	age := c.Param("age")
	c.JSON(200, gin.H{
		"name": name,
		"age":  age,
	})
}

func main() {
	// NOTE: See we’re using = to assign the global var
	// instead of := which would assign it only in this function
	//db, err = gorm.Open("sqlite3", "./gorm.db")
	db, err = gorm.Open("postgres", "host=taco-db port=5432 user=taco dbname=taco password=pass1234 sslmode=disable")
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()
	db.AutoMigrate(&handler.Person{})
	db.AutoMigrate(&stock.Product{})
	db.AutoMigrate(&stock.Price{})
	r := gin.Default()

	// gredis.init()

	// template test
	// r.GET("/people/", handler.GetPeople(db))

	// person: orm test
	r.GET("/people/", handler.GetPeople(db))
	r.GET("/people/:id", handler.GetPerson(db))
	r.POST("/people", handler.CreatePerson(db))
	r.PUT("/people/:id", handler.UpdatePerson(db))
	r.DELETE("/people/:id", handler.DeletePerson(db))

	// Newsfeed test
	feed := newsfeed.New()

	fmt.Println("Hello World")

	// r := gin.Default()
	r.GET("/", HomePage)
	r.POST("/", PostHomePage)
	r.GET("/query", QueryStrings) // /query?name=ellen&age=24
	r.GET("/path/:name/:age", PathParameters)

	r.GET("/ping", handler.PingGet())
	r.GET("/newsfeed", handler.NewsfeedGet(feed))
	r.POST("/newsfeed", handler.NewsfeedPost(feed))

	r.POST("/BrokerSend", handler.BrokerSend())
	r.POST("/BrokerReceive", handler.BrokerReceive())

	r.POST("/broadcast/:message", handler.BroadcastMessage())
	r.POST("/consumer/*username", handler.Consumer())

	// r.GET("/pathtest/:name", func(){
	// 	name
	// })

	r.Run()
}
