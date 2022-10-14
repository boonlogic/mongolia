package main

import (
	"encoding/json"
	"gitlab.boonlogic.com/development/expert/mongolia/pkg/odm"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

	//preValidate := func(any) *odm.Model {
	//	fmt.Println("prevalidating...")
	//	return nil
	//}
	//hooks := &odm.Hooks{
	//	PreValidate: preValidate,
	//}

	//roles := odm.GetModel("roles")
	//doc, err := roles.CreateOne(role)
	//if err != nil {
	//	log.Fatalf("failed to CreateOne: %s\n", err)
	//}
	//fmt.Printf("created document:\n%+v\n", doc)

	type Role struct {
		Name        string               `json:"name"`
		Permissions []primitive.ObjectID `json:"permissions"`
	}
	role := &Role{
		Name:        "brads-role",
		Permissions: []primitive.ObjectID{},
	}

	roleText, err := json.Marshal(role)
	if err != nil {
		log.Fatalf("failed to marshal struct to json")
	}

	path := "schemas/role.json"
	schemaText, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("failed to load role schema: %s\n", err)
	}

	ok, err := odm.ValidateJSONSchema(roleText, schemaText)
	if err != nil {
		log.Fatalf("failed to run validator: %s\n", err)
	}

	if !ok {
		log.Fatalln("invalid model")
	}
	log.Println("model is valid")

	//odm.RegisterModel("roles", schemaText, hooks)

	//router := gin.Default()
	//router.GET("/hello", controllers.SayHello)
	//router.Run("localhost:8080")
}
