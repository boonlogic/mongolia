package main

import (
	"fmt"
	"gitlab.boonlogic.com/development/expert/mongolia/pkg/odm"
	"io/ioutil"
	"log"
)

func main() {
	if err := odm.Configure(); err != nil {
		log.Fatalf("failed to configure: %s\n", err)
	}

	if err := odm.Connect(); err != nil {
		log.Fatalf("failed to connect: %s\n", err)
	}

	odm.Drop()

	path := "schemas/role.json"
	schemaText, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("failed to load role schema: %s\n", err)
	}

	preValidate := func(any) *odm.Model {
		fmt.Println("prevalidating...")
		return nil
	}
	hooks := &odm.Hooks{
		PreValidate: preValidate,
	}
	if err := odm.RegisterModel("roles", schemaText, hooks); err != nil {
		log.Fatalf("failed to register model: %s\n", err)
	}

	type Role struct {
		ID          string   `json:"id"`
		Name        string   `json:"name"`
		Permissions []string `json:"permissions"`
	}
	r := &Role{
		ID:          "6349a84fe97051c7b555e172",
		Name:        "admin",
		Permissions: []string{"+:*:*"},
	}

	roles := odm.GetModel("roles")
	doc, err := roles.CreateOne(r)
	if err != nil {
		log.Fatalf("failed to CreateOne: %s\n", err)
	}
	fmt.Printf("created document:\n%+v\n", doc)



	//router := gin.Default()
	//router.GET("/hello", controllers.SayHello)
	//router.Run("localhost:8080")
}
