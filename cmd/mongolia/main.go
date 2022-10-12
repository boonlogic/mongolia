package main

import (
	"github.com/gin-gonic/gin"
	"gitlab.boonlogic.com/development/expert/mongolia/pkg/controllers"
)

func main() {
	router := gin.Default()
	router.GET("/hello", controllers.SayHello)
	router.Run("localhost:8080")
}
