package mongolia

import (
	"errors"
	"fmt"
	"io/ioutil"
)

type ODM struct {
	compiler *Compiler
	schemas  map[string]Schema
	colls    map[string]Collection
}

func NewODM() *ODM {
	return &ODM{
		compiler: NewCompiler(),
		schemas:  make(map[string]Schema),
		colls:    make(map[string]Collection),
	}
}

func (o *ODM) AddSchema(name string, path string) error {
	// If AddResource succeeds, the spec is valid jsonschema.
	def, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	if err := odm.compiler.AddResource(path, def); err != nil {
		return err
	}
	schema, err := odm.compiler.Compile(path)
	if err != nil {
		return err
	}
	odm.schemas[name] = *schema
	return nil
}

func (o *ODM) GetCollection(name string) (*Collection, error) {
	coll, ok := o.colls[name]
	if !ok {
		return nil, errors.New(fmt.Sprintf("no collection named \"%s\"", name))
	}
	return &coll, nil
}
