package main

import (
	"fmt"

	"github.com/boonlogic/mongolia"
)

func main() {
	odm := mongolia.NewODM()

	fmt.Println(odm.URI)
}
