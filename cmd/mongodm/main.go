package main

import (
	"github.com/gin-gonic/gin"
	"gitlab.boonlogic.com/development/expert/mongolia/pkg/mongodm"
	"gitlab.boonlogic.com/development/expert/mongolia/pkg/restapi"
	"log"
)

func main() {
	if err := mongodm.Configure(); err != nil {
		log.Fatalf("failed to configure: %s\n", err)
	}
	if err := mongodm.Connect(); err != nil {
		log.Fatalf("failed to connect: %s\n", err)
	}
	router := gin.Default()
	router.GET("/hello", restapi.SayHello)
	if err := router.Run("localhost:8080"); err != nil {
		log.Fatalf("error: %s\n", err)
	}
}
