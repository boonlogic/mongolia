package main

import (
	"github.com/gin-gonic/gin"
	"gitlab.boonlogic.com/development/expert/mongolia/pkg/controllers"
	"gitlab.boonlogic.com/development/expert/mongolia/pkg/odm"
	"log"
)

func main() {
	if err := odm.Configure(); err != nil {
		log.Fatalf("failed to configure: %s\n", err)
	}

	if err := odm.Connect(); err != nil {
		log.Fatalf("failed to connect: %s\n", err)
	}

	//perm := &odm.Schema{
	//	Name: "permission",
	//	Attributes: []odm.Attribute{
	//		{
	//			Name:     "name",
	//			Unique:   true,
	//			Required: true,
	//			Type:     types.STRING,
	//		},
	//		{
	//			Name:     "token",
	//			Unique:   true,
	//			Required: true,
	//			Type:     types.STRING,
	//			Pattern:  "^(\\+|\\-):([a-zA-Z0-9]|\\*):([a-zA-Z0-9]|\\*)$",
	//		},
	//		{
	//			Name:     "access",
	//			Required: true,
	//			Type:     types.STRING,
	//			Enum: odm.Enum{
	//				Values: []Value{}{
	//					"allow",
	//					"deny",
	//				},
	//			},
	//		},
	//		{
	//			Name:     "resource",
	//			Type:     types.STRING,
	//			Required: true,
	//			Pattern:  "^([a-zA-Z0-9-_]|\\*)$",
	//		},
	//		{
	//			Name:     "action",
	//			Type:     types.STRING,
	//			Required: true,
	//			Pattern:  "^([a-zA-Z0-9-_]|\\*)$",
	//		},
	//	},
	//}
	//
	//role := &odm.Schema{
	//	Name: "role",
	//	Attributes: []odm.Attribute{
	//		{
	//			Name:     "name",
	//			Type:     types.STRING,
	//			Unique:   true,
	//			Required: true,
	//			Pattern:  "^[a-zA-Z0-9-_]*$",
	//		},
	//		{
	//			Name: "permissions",
	//			Type: types.ARRAY,
	//			Required: true,
	//			Item: odm.ArrayItem{
	//				Type: types.POINTER,
	//				Collection: "permissions",
	//			},
	//		},
	//		{
	//			Name:      "user emails",
	//			Type:      types.ARRAY,
	//			Item: odm.ArrayItem{
	//				Type: types.STRING,
	//				Format: "email",
	//			},
	//		},
	//	},
	//}

	//if err := odm.AddSchema(perm); err != nil {
	//	panic(any(err))
	//}
	//if err := odm.AddSchema(role); err != nil {
	//	panic(any(err))
	//}

	// add the rest of the schemas: perms and roles (finish defining the spec)
	// draw solid dotted line, begin using the spec
	// implement odm.Query, odm.Find etc and try them

	router := gin.Default()
	router.GET("/hello", controllers.SayHello)
	router.Run("localhost:8080")
}
