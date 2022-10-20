package mongodm

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io/ioutil"
	"testing"
)

func TestSmoke(t *testing.T) {
	cfg := NewConfigFromEnvironment()
	err := Connect(*cfg)
	require.Nil(t, err)
	drop()

	// Load schema definition.
	path := "schemas/role.json"
	def, err := ioutil.ReadFile(path)
	require.Nil(t, err)

	// Define any hook functions for this schema.
	hooks := myHooks()

	// Register our schema + hooks as an ODM model, getting a Collection
	// which is a handle to the object store that returns Documents.
	AddSchema("roles", def, hooks)
	roles, ok := GetCollection("roles")
	require.True(t, ok)

	// Define the struct which will be mapped to/from the model.
	type Role struct {
		ID          primitive.ObjectID `json:"id" bson:"id"`
		Name        string             `json:"name" bson:"name"`
		Permissions []string           `json:"permissions" bson:"permissions"`
	}

	// Create a role document from the Role struct.
	r := &Role{
		Name:        "admin",
		Permissions: []string{"+:*:*"},
	}
	doc, err := roles.CreateOne(r)
	require.Nil(t, err)
	require.NotNil(t, doc)

	// Creation should fail if username does not match pattern regex.
	rbad := *r
	rbad.Name = "b@dusern@me"
	doc, err = roles.CreateOne(&rbad)
	require.Nil(t, doc)
	require.NotNil(t, err)

	// Find a previously inserted document.
	doc, err = roles.FindOne(nil)
	require.Nil(t, err)
	require.NotNil(t, doc)

	// Create several documents.
	obj1 := *r
	obj1.ID = primitive.NewObjectID()
	obj2 := *r
	obj2.ID = primitive.NewObjectID()
	obj3 := *r
	obj3.ID = primitive.NewObjectID()
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

	// Update many documents.
	docs, err = roles.UpdateMany(nil)
	require.Nil(t, err)
	require.Len(t, docs, 4)

	// Delete one document.
	doc, err = roles.RemoveOne(nil)
	require.Nil(t, err)
	require.NotNil(t, doc)

	// Delete many documents.
	docs, err = roles.RemoveMany(nil)
	require.Nil(t, err)
	require.Len(t, docs, 3)
}

func myHooks() *Hooks {
	preValidate := func(*Document) error {
		fmt.Println("hello from preValidate")
		return nil
	}
	postValidate := func(*Document) error {
		fmt.Println("hello from postValidate")
		return nil
	}
	preCreate := func(*Document) error {
		fmt.Println("hello from preCreate")
		return nil
	}
	preUpdate := func(*Document) error {
		fmt.Println("hello from preUpdate")
		return nil
	}
	preSave := func(*Document) error {
		fmt.Println("hello from preSave")
		return nil
	}
	preRemove := func(*Document) error {
		fmt.Println("hello from preRemove")
		return nil
	}
	postCreate := func(*Document) error {
		fmt.Println("hello from postCreate")
		return nil
	}
	postUpdate := func(*Document) error {
		fmt.Println("hello from postUpdate")
		return nil
	}
	postSave := func(*Document) error {
		fmt.Println("hello from postSave")
		return nil
	}
	postRemove := func(*Document) error {
		fmt.Println("hello from postRemove")
		return nil
	}
	return &Hooks{
		PreValidate:  preValidate,
		PostValidate: postValidate,
		PreSave:      preSave,
		PostSave:     postSave,
		PreCreate:    preCreate,
		PostCreate:   postCreate,
		PreUpdate:    preUpdate,
		PostUpdate:   postUpdate,
		PreRemove:    preRemove,
		PostRemove:   postRemove,
	}
}
