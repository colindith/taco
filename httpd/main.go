package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"taco/httpd/handler"
	"taco/httpd/handler/stock"
	"taco/packages/gredis"
	"taco/platform/newsfeed"
	"time"

	"github.com/dgrijalva/jwt-go"
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

	{
		r.POST("/login", login)
		// r.GET("/logout", logout)
		private := r.Group("/private")
		private.Use(AuthRequired())
		private.POST("/auth_test/:token", func(c *gin.Context) {
			c.JSON(http.StatusOK, "get private data~~~!!")
		})
		private.POST("/head_auth/", func(c *gin.Context) {
			c.JSON(http.StatusOK, "get head private data~~~!!")
		})
	}

	r.Run()
}

const (
	// make this const
	mySigningKey = "WOW,MuchShibe,ToDogge"
)

func respondWithError(code int, message string, c *gin.Context) {
	resp := map[string]string{"error": message}
	c.JSON(code, resp)
	c.Abort()
}

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		authorization, exists := c.Request.Header["Authorization"]

		if exists != true {
			respondWithError(401, "Error in parsing: "+err.Error(), c)
			return
		}
		is_authenticated := parseJWT(strings.Split(authorization[0], " ")[1])

		if is_authenticated == false {
			// c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session token"})
			respondWithError(401, "Invalid Authorization", c)
			return
		} else {
			c.Next()
		}
	}
}

// curl -X POST -F 'username=hello' -F 'password=itsme' 'http://localhost:8080/login'
// curl -X POST -H "Authorization: JWT $token" 'http://localhost:8080/private/head_auth/'

func login(c *gin.Context) {
	// session := sessions.Default(c)
	username := c.PostForm("username")
	password := c.PostForm("password")

	if strings.Trim(username, " ") == "" || strings.Trim(password, " ") == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Parameters can't be empty"})
		return
	}
	// check if username in db
	if username == "hello" && password == "itsme" {
		// session.Set("user", username) //In real world usage you'd set this to the users ID
		// err := session.Save()
		token, err := newJWT(username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate JWT token: " + err.Error()})
		} else {
			c.JSON(http.StatusOK, gin.H{"token": token})
		}
	}
}

func newJWT(username string) (string, error) {
	// Create the token
	token := jwt.New(jwt.SigningMethodHS256)
	// Set some claims
	claims := make(jwt.MapClaims)
	claims["username"] = "username"
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
	token.Claims = claims
	// Sign and get the complete encoded token as a string
	// TODO: here should use byte array directly
	tokenString, err := token.SignedString([]byte(mySigningKey))
	return tokenString, err
}

func parseJWT(myToken string) bool {
	token, err := jwt.Parse(myToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(mySigningKey), nil
	})
	if err != nil {
		fmt.Println("Error occured.  parseJWT" + err.Error())
		return false
	} else if token.Valid {
		fmt.Println("Your token is valid.  I like your style.")
		return true
	} else {
		fmt.Println("This token is terrible!  I cannot accept this.")
		return false
	}
}
