package main

import (
	"github.com/gin-gonic/gin"
	"gitlab.boonlogic.com/development/expert/mongolia/mongolia"
	"gitlab.boonlogic.com/development/expert/mongolia/restapi"
	"log"
	"os"
	"strconv"
	"time"
)

func main() {
	n, err := strconv.Atoi(os.Getenv("MONGOLIA_TIMEOUT"))
	if err != nil {
		log.Fatalf("could not convert MONGOLIA_TIMEOUT to int: '%s'\n", os.Getenv("MONGOLIA_TIMEOUT"))
	}
	timeout := time.Duration(n) * time.Second

	cfg := mongolia.DefaultConfig().
		SetURI(os.Getenv("MONGOLIA_URI")).
		SetDBName(os.Getenv("MONGOLIA_DB_NAME")).
		SetTimeout(timeout)

	if err := mongolia.Connect(cfg); err != nil {
		log.Fatalf("failed to connect: %s\n", err)
	}

	router := gin.Default()
	router.GET("/hello", restapi.SayHello)
	if err := router.Run("localhost:8080"); err != nil {
		log.Fatalf("error: %s\n", err)
	}
}
