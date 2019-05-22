package handler

import (
	"net/http"
	"taco/packages/broker"

	"github.com/gin-gonic/gin"
)

func BrokerSend() gin.HandlerFunc {
	return func(c *gin.Context) {
		broker.Send()
		c.Status(http.StatusNoContent)
	}
}

func BrokerReceive() gin.HandlerFunc {
	return func(c *gin.Context) {
		broker.Receive()
		c.Status(http.StatusNoContent)
	}
}

func BroadcastMessage() gin.HandlerFunc {
	return func(c *gin.Context) {
		msg := c.Param("message")
		broker.BroadcastMessage(msg)
		c.Status(http.StatusNoContent)
	}
}

func Consumer() gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.Param("username")
		broker.Consumer(username)
		c.Status(http.StatusNoContent)
	}
}
