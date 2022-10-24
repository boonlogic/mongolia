package mongodm

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"gitlab.boonlogic.com/development/expert/mongolia/pkg/mongodm/options"
	"gitlab.boonlogic.com/development/expert/mongolia/pkg/mongodm/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"io/ioutil"
	"testing"
)

func TestSmoke(t *testing.T) {

	mongo.Client
	// Setup
	cfg := options.NewConfigFromEnvironment()
	err := Connect(cfg)
	require.Nil(t, err)

	options.drop() // cleanup from last test

	spec1, err := types.NewSpecFromFile("schemas/role.json")
	require.Nil(t, err)
	data, err := ioutil.ReadFile("schemas/role.json")
	require.Nil(t, err)
	spec2, err := types.NewSpecFromJSON(data)
	require.Nil(t, err)

	hooks1 := myHooks()
	hooks2 := myHooks()

	err = AddSchema("roles", *spec1, hooks1)
	require.Nil(t, err)
	err = AddSchema("roles2", *spec2, hooks2)
	require.Nil(t, err)

	roles, ok := GetCollection("roles")
	require.True(t, ok)
	_, ok = GetCollection("roles2")
	require.True(t, ok)

	// Usage
	type Role struct {
		ID          primitive.ObjectID `json:"id" bson:"id"`
		Name        string             `json:"name" bson:"name"`
		Permissions []string           `json:"permissions" bson:"permissions"`
	}
	r := &Role{
		Name:        "admin",
		Permissions: []string{"+:*:*"},
	}

	// Create a role document from the Role struct.
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

func myHooks() *types.Hooks {
	preValidate := func(*types.Document) error {
		fmt.Println("hello from preValidate")
		return nil
	}
	postValidate := func(*types.Document) error {
		fmt.Println("hello from postValidate")
		return nil
	}
	preCreate := func(*types.Document) error {
		fmt.Println("hello from preCreate")
		return nil
	}
	preUpdate := func(*types.Document) error {
		fmt.Println("hello from preUpdate")
		return nil
	}
	preSave := func(*types.Document) error {
		fmt.Println("hello from preSave")
		return nil
	}
	preRemove := func(*types.Document) error {
		fmt.Println("hello from preRemove")
		return nil
	}
	postCreate := func(*types.Document) error {
		fmt.Println("hello from postCreate")
		return nil
	}
	postUpdate := func(*types.Document) error {
		fmt.Println("hello from postUpdate")
		return nil
	}
	postSave := func(*types.Document) error {
		fmt.Println("hello from postSave")
		return nil
	}
	postRemove := func(*types.Document) error {
		fmt.Println("hello from postRemove")
		return nil
	}
	return &types.Hooks{
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
