package v0

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/require"
	"gitlab.boonlogic.com/development/expert/mongolia/mongolia/v0/options"
	"io/ioutil"
	"testing"
)

type Role struct {
	Name        string   `json:"name"  bson:"name"`
	Permissions []string `json:"permissions" bson:"permissions"`
}

func TestSmoke(t *testing.T) {
	opts := options.ODM().
		SetHost("localhost").
		SetName("mongolia-local").
		SetEphemeral(true)
	opts.Validate()

	err := Connect(opts)
	require.Nil(t, err)

	Drop()

	buf, err := ioutil.ReadFile("schemas/role.json")
	if err != nil {
		panic(err)
	}
	spec := NewSpec(buf)
	err = spec.Validate()
	require.Nil(t, err)

	spec2, err := NewSpecFromFile("schemas/role.json")
	require.Nil(t, err)
	err = spec2.Validate()
	require.Nil(t, err)

	err = AddSchema("roles", spec, myHooks())
	require.Nil(t, err)

	coll, err := GetCollection("roles")
	require.Nil(t, err)

	role := &Role{
		Name:        "admin",
		Permissions: []string{"+:*:*"},
	}

	var doc map[string]any
	buf, err = json.Marshal(role)
	if err != nil {
		panic(err)
	}
	if err = json.Unmarshal(buf, &doc); err != nil {
		panic(err)
	}
	model, err := coll.CreateOne(doc)
	require.Nil(t, err)
	require.NotNil(t, model)

	fmt.Printf("model: %+v\n", model)

	//// Creation should fail if username does not match pattern regex.
	//rbad := *r
	//rbad.Name = "b@dusern@me"
	//model, err = roles.CreateOne(&rbad)
	//require.Nil(t, model)
	//require.NotNil(t, err)
	//
	//// Find a previously inserted document.
	//model, err = roles.FindOne(nil)
	//require.Nil(t, err)
	//require.NotNil(t, model)
	//
	//// Create several documents.
	//obj1 := *r
	//obj1.ID = primitive.NewObjectID()
	//obj2 := *r
	//obj2.ID = primitive.NewObjectID()
	//obj3 := *r
	//obj3.ID = primitive.NewObjectID()
	//objs := []any{obj1, obj2, obj3}
	//docs, err := roles.CreateMany(objs)
	//require.Nil(t, err)
	//require.NotNil(t, docs)
	//
	//// Find all previously created documents.
	//docs, err = roles.FindMany(nil)
	//require.Nil(t, err)
	//require.Len(t, docs, 4)
	//
	//// Update a document.
	//model, err = roles.UpdateOne(nil)
	//require.Nil(t, err)
	//require.NotNil(t, model)
	//
	//// Update many documents.
	//docs, err = roles.UpdateMany(nil)
	//require.Nil(t, err)
	//require.Len(t, docs, 4)
	//
	//// Delete one document.
	//model, err = roles.RemoveOne(nil)
	//require.Nil(t, err)
	//require.NotNil(t, model)
	//
	//// Delete many documents.
	//docs, err = roles.RemoveMany(nil)
	//require.Nil(t, err)
	//require.Len(t, docs, 3)
}

func myHooks() *Hooks {
	preValidate := func(*Model) error {
		fmt.Println("hello from preValidate")
		return nil
	}
	postValidate := func(*Model) error {
		fmt.Println("hello from postValidate")
		return nil
	}
	preCreate := func(*Model) error {
		fmt.Println("hello from preCreate")
		return nil
	}
	preUpdate := func(*Model) error {
		fmt.Println("hello from preUpdate")
		return nil
	}
	preSave := func(*Model) error {
		fmt.Println("hello from preSave")
		return nil
	}
	preRemove := func(*Model) error {
		fmt.Println("hello from preRemove")
		return nil
	}
	postCreate := func(*Model) error {
		fmt.Println("hello from postCreate")
		return nil
	}
	postUpdate := func(*Model) error {
		fmt.Println("hello from postUpdate")
		return nil
	}
	postSave := func(*Model) error {
		fmt.Println("hello from postSave")
		return nil
	}
	postRemove := func(*Model) error {
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
