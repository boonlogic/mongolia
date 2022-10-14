package odm

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
)


var modelRegistry = ModelRegistry{}

type ModelRegistry map[string]Schema

func GetModel(name string) Schema {
	return modelRegistry[name]
}

type Model interface {
	PreValidate(any) *Model
	PostValidate(any) *Model
	PreSave(any) *Model
	PostSave(any) *Model
	PreCreate(any) *Model
	PostCreate(any) *Model
	PreUpdate(any) *Model
	PostUpdate(any) *Model
	PreRemove(any) *Model
	PostRemove(any) *Model

	////////////////////////////////

	CreateOne(obj any) (doc *Document, err error)
	CreateMany(any) ([]Document, error)
	FindOne(any) (*Document, error)
	FindMany(any) ([]Document, error)
	UpdateOne(any) (*Document, error)
	UpdateMany(any) ([]Document, error)
	RemoveOne(any) (*Document, error)
	RemoveMany(any) ([]Document, error)
}

func RegisterModel(name string, spec []byte, hooks *Hooks) {
	// todo: validate that it is good OpenAPI
	//loader := &openapi3.Loader{Context: ctx(), IsExternalRefsAllowed: true}
	//doc, _ := loader.LoadFromFile(".../My-OpenAPIv3-API.yml")
	//_ := doc.Validate(ctx())

	s := Schema{
		Name: name,
		Definition: spec,
		Hooks:      hooks,
	}
	modelRegistry[name] = s
}

type Hooks struct {
	PreValidate  func(any) *Model
	PostValidate func(any) *Model
	PreSave      func(any) *Model
	PostSave     func(any) *Model
	PreCreate    func(any) *Model
	PostCreate   func(any) *Model
	PreUpdate    func(any) *Model
	PostUpdate   func(any) *Model
	PreRemove    func(any) *Model
	PostRemove   func(any) *Model
}

type Document bson.D

func ctx() context.Context {
	return context.Background()
}
