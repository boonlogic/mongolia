package main

import (
	"fmt"

	"gitlab.boonlogic.com/development/expert/mongolia/mongolia"
)

func main() {
	odm := mongolia.NewODM()

	fmt.Println(odm.URI)
}
