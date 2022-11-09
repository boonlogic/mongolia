package mongolia

import (
	"context"
)

var odm *ODM

func init() {
	odm = NewODM()
}

func ctx() context.Context {
	return context.Background()
}

func AddSchema(name string, path string) error {
	return odm.AddSchema(name, path)
}

func GetCollection(name string) (*Collection, error) {
	return odm.GetCollection(name)
}
