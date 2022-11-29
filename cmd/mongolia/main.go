package main

import (
	"github.com/gin-gonic/gin"
	"gitlab.boonlogic.com/development/expert/mongolia/mongolia"
	"gitlab.boonlogic.com/development/expert/mongolia/restapi"
	"log"
	"os"
)

func main() {
	cfg := mongolia.DefaultConfig().
		SetURI(os.Getenv("MONGOLIA_URI")).
		SetDBName(os.Getenv("MONGOLIA_DBNAME"))

	if err := mongolia.Connect(cfg); err != nil {
		log.Fatalf("failed to connect: %s\n", err)
	}

	router := gin.Default()
	router.GET("/hello", restapi.SayHello)
	if err := router.Run("localhost:8080"); err != nil {
		log.Fatalf("error: %s\n", err)
	}
}
