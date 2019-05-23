package user

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	// ID       uint   `json:"id" gorm:"AUTO_INCREMENT;PRIMARY_KEY"`
	Username string `json:"username" gorm:"unique;not null"`
	Password string `json:"password"`
}

func Registry(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user User
		c.BindJSON(&user)
		fmt.Println(user.Username)
		fmt.Println(user.Password)
		if dbc := db.Create(&user); dbc.Error != nil {
			c.JSON(200, gin.H{
				"success": false,
				"msg":     dbc.Error,
			})
			// TODO: dbc.Error is a JSON like object, parse the message in
		}
		c.JSON(200, gin.H{
			"success":  true,
			"username": user.Username,
		})
	}
}

func ListUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var users []User
		if err := db.Find(&users).Error; err != nil {
			c.AbortWithStatus(404)
			fmt.Println(err)
		} else {
			c.JSON(200, users)
		}
	}
}

func Login(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Here realized the post data fetching with out using the go bindind
		byte_data, err := c.GetRawData()
		if err != nil {

		}
		// fmt.Println(byte_data)
		data := make(map[string]string)
		json.Unmarshal(byte_data, &data)

		username := data["username"]
		password := data["password"]
		// fmt.Println("----------> ", username, password)

		if strings.Trim(username, " ") == "" || strings.Trim(password, " ") == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Parameters can't be empty"})
			return
		}
		var user User
		if result := db.Where("Username=(?)", username).First(&user); result.Error != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "No such username"})
			return
		}
		// fmt.Print("user found: ", user)
		if user.Password != password {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Wrong password"})
			return
		}

		token, err := newJWT(username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate JWT token: " + err.Error()})
			return
		} else {
			c.JSON(http.StatusOK, gin.H{"token": token})
		}
	}
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
			respondWithError(401, "No Authorization: ", c)
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

// curl -X POST -d '{"username":"peter2","password":"123456"}' 'http://localhost:8080/user/login'
// curl -X POST -H "Authorization: JWT $token" 'http://localhost:8080/private/head_auth/'

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
		fmt.Println("Error occured when parseJWT: " + err.Error())
		return false
	} else if token.Valid {
		return true
	} else {
		return false
	}
}
