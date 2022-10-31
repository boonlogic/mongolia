package main

import (
	"github.com/gin-gonic/gin"
	"gitlab.boonlogic.com/development/expert/mongolia/mongodm/v0"
	"gitlab.boonlogic.com/development/expert/mongolia/mongodm/v0/options"
	"gitlab.boonlogic.com/development/expert/mongolia/restapi"
	"log"
)

func main() {
	opts := options.ODM().
		SetCloud(false).
		SetHost("localhost").
		SetName("mongodm-local").
		SetEphemeral(false)

	if err := v0.Connect(opts); err != nil {
		log.Fatalf("failed to configure: %s\n", err)
	}

	router := gin.Default()
	router.GET("/hello", restapi.SayHello)
	if err := router.Run("localhost:8080"); err != nil {
		log.Fatalf("error: %s\n", err)
	}
}
