package odm

import (
	"net/url"
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

	CreateOne(obj any) (doc *Document, err error)
	CreateMany(any) ([]Document, error)
	FindOne(any) (*Document, error)
	FindMany(any) ([]Document, error)
	UpdateOne(any) (*Document, error)
	UpdateMany(any) ([]Document, error)
	RemoveOne(any) (*Document, error)
	RemoveMany(any) ([]Document, error)
}

func RegisterModel(name string, spec []byte, hooks *Hooks) error {
	vfunc, err := ValidateSpec(name, spec)
	if err != nil {
		return err
	}
	s := Schema{
		Name:       name,
		Definition: spec,
		Validator:  vfunc,
		Hooks:      hooks,
	}
	modelRegistry[name] = s
	return nil
}

func encodeURL(s string) string {
	return url.QueryEscape(s)
}
