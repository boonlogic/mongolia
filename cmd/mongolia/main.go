package main

import (
	"github.com/gin-gonic/gin"
	"gitlab.boonlogic.com/development/expert/mongolia/pkg/restapi"
	"gitlab.boonlogic.com/development/expert/mongolia/pkg/odm"
	"log"
)

func main() {
	if err := odm.Configure(); err != nil {
		log.Fatalf("failed to configure: %s\n", err)
	}
	if err := odm.Connect(); err != nil {
		log.Fatalf("failed to connect: %s\n", err)
	}
	router := gin.Default()
	router.GET("/hello", restapi.SayHello)
	router.Run("localhost:8080")
}
