package main

import (
	"log"
	"os"

	"gitlab.boonlogic.com/development/expert/mongolia/mongolia"
)

func main() {
	cfg := mongolia.DefaultConfig().
		SetURI(os.Getenv("MONGOLIA_URI")).
		SetDBName(os.Getenv("MONGOLIA_DBNAME"))

	if err := mongolia.Connect(cfg); err != nil {
		log.Fatalf("failed to connect: %s\n", err)
	}

}
