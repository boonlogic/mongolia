package odm

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"testing"
)

func TestSmoke(t *testing.T) {
	err := Configure()
	require.Nil(t, err)
	err = Connect()
	require.Nil(t, err)
	err = Drop()
	require.Nil(t, err)

	// Load test schema.
	path := "schemas/role.json"
	schemaText, err := ioutil.ReadFile(path)
	require.Nil(t, err)

	// Create an example hook function.
	preValidate := func(any) *Model {
		fmt.Println("prevalidating...")
		return nil
	}
	hooks := &Hooks{
		PreValidate: preValidate,
	}

	// Register our schema as a model with the given hooks.
	err = RegisterModel("roles", schemaText, hooks)
	require.Nil(t, err)

	// Define the struct which will be mapped to/from the model.
	type Role struct {
		ID          string   `json:"id" bson:"id"`
		Name        string   `json:"name" bson:"name"`
		Permissions []string `json:"permissions" bson:"permissions"`
	}

	// Get the roles Model, which is an interface to the roles collection.
	roles := GetModel("roles")

	// Create a role document from the Role struct.
	r := &Role{
		ID:          "6349a84fe97051c7b555e172",
		Name:        "admin",
		Permissions: []string{"+:*:*"},
	}
	doc, err := roles.CreateOne(r)
	require.Nil(t, err)
	require.NotNil(t, doc)

	// Creation should fail if username does not match pattern regex.
	rbad := *r
	rbad.Name = "bad@#$"
	doc, err = roles.CreateOne(&rbad)
	require.Nil(t, doc)
	require.NotNil(t, err)

	// Find a previously inserted document.
	doc, err = roles.FindOne(nil)
	require.Nil(t, err)
	require.NotNil(t, doc)

	// Create several documents.
	obj1 := *r
	obj1.ID = "6349a84fe97051c7b555e173"
	obj2 := *r
	obj2.ID = "6349a84fe97051c7b555e174"
	obj3 := *r
	obj3.ID = "6349a84fe97051c7b555e175"
	objs := []any{obj1, obj2, obj3}
	docs, err := roles.CreateMany(objs)
	require.Nil(t, err)
	require.NotNil(t, docs)

	// Find all previously created documents.
	docs, err = roles.FindMany(nil)
	require.Nil(t, err)
	require.Len(t, docs, 4)

	// Update a document.
	doc, err = roles.UpdateOne(nil)
	require.Nil(t, err)
	require.NotNil(t, doc)
}