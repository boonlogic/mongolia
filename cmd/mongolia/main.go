package main

import (
	"fmt"

	"github.com/boonlogic/mongolia/mongolia"
)

func main() {
	odm := mongolia.NewODM()

	fmt.Println(odm.URI)
}
