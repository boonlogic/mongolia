package main

import (
	"github.com/gin-gonic/gin"
	"gitlab.boonlogic.com/development/expert/mongolia/mongodm"
	"gitlab.boonlogic.com/development/expert/mongolia/mongodm/options"
	"gitlab.boonlogic.com/development/expert/mongolia/restapi"
	"log"
)

func main() {
	opts := options.Configure().
		SetCloud(false).
		SetHost("localhost").
		SetDatabase("mongodm-local").
		SetEnvironment("dev")

	if err := mongodm.Configure(opts); err != nil {
		log.Fatalf("failed to configure: %s\n", err)
	}

	router := gin.Default()
	router.GET("/hello", restapi.SayHello)
	if err := router.Run("localhost:8080"); err != nil {
		log.Fatalf("error: %s\n", err)
	}
}
