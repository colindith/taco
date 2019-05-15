package handler

import (
	"fmt"
	"net/http"
	"taco/platform/newsfeed"
	"github.com/gin-gonic/gin"
)

type newsfeedPostRequest struct {
	Title string `json:"title"`
	Post 	string `json:"post"`
}

func NewsfeedPost(feed newsfeed.Adder) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestBody := newsfeedPostRequest{}
		c.Bind(&requestBody)

		fmt.Println("requestBody", requestBody)


		item := newsfeed.Item{
			Title: requestBody.Title,
			Post: requestBody.Post,
		}
		feed.Add(item)

		c.Status(http.StatusNoContent)
	}
}