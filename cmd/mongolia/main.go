package main

import (
	"fmt"
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

	// add the rest of the schemas: perms and roles (finish defining the spec)
	path := "/Users/lukearend/builder/packages/mongolia/mongolia/schemas/role.json"
	spec, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("failed to load license schema: %s\n", err)
	}

	preValidate := func(any) *odm.Model {
		fmt.Println("prevalidating...")
		return nil
	}
	hooks := &odm.Hooks{
		PreValidate: preValidate,
	}

	odm.RegisterModel("roles", spec, hooks)

	//role := &Role{
	//	Name:        "brads-first-role",
	//	Permissions: []primitive.ObjectID{},
	//}
	//roles := odm.GetModel("roles")
	//doc, err := roles.CreateOne(role)
	//if err != nil {
	//	log.Fatalf("failed to CreateOne: %s\n", err)
	//}
	//
	//fmt.Printf("created document:\n%+v\n", doc)

	schema := `
{
  "$id": "https://gitlab.com/boonlogic/development/expert-api/api/amberv2/schemas/role.json",
  "title": "Role",
  "type": "object",
  "required": [ "productId", "productName", "price" ]
  "properties": {
    "_id": {
      "type": "objectId"
    },
    "name": {
      "type": "string",
      "pattern"": "^[a-zA-Z0-9-_]*$"
    },
    "price": {
      "description": "The price of the product",
      "type": "number",
      "exclusiveMinimum": 0
    },
    "tags": {
      "description": "Tags for the product",
      "type": "array",
      "items": {
        "type": "string"
      },
      "minItems": 1,
      "uniqueItems": true
    }
  },
}
`

	role := map[string]interface{}{
		"name": "brads-role",
		"permissions": []primitive.ObjectID{},
	}

	odm.Validate()

	//router := gin.Default()
	//router.GET("/hello", controllers.SayHello)
	//router.Run("localhost:8080")
}
