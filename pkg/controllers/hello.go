package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func SayHello(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, gin.H{"message": "hello client!"})
}
