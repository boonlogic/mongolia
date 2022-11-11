package mongolia

import (
	"errors"
	"fmt"
	"io/ioutil"
)

type ODM struct {
	config   *Config
	compiler *Compiler
	colls    map[string]Collection
}

func NewODM() *ODM {
	return &ODM{
		compiler: NewCompiler(),
		colls:    make(map[string]Collection),
	}
}

func (o *ODM) Connect(config *Config) error {
	if config == nil {
		config = NewConfig()
	} else {
		config = DefaultConfig().Merge(config)
	}
	if err := connect(config); err != nil {
		return err
	}
	return nil
}

func (o *ODM) AddSchema(name string, path string) error {
	def, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	schema, err := odm.compileSchema(path, def)
	if err != nil {
		return err
	}
	coll, err := connectCollection(name, schema)
	if err != nil {
		return nil
	}
	odm.colls[name] = *coll
	return nil
}

func (o *ODM) compileSchema(path string, def []byte) (*Schema, error) {
	r := &Resource{
		URI:        path,
		Definition: def,
	}
	if err := odm.compiler.AddResource(path, r); err != nil {
		return nil, err
	}
	schema, err := odm.compiler.Compile(path)
	if err != nil {
		return nil, err
	}
	return schema, nil
}

func (o *ODM) GetCollection(name string) (*Collection, error) {
	coll, ok := o.colls[name]
	if !ok {
		return nil, errors.New(fmt.Sprintf("no collection named \"%s\"", name))
	}
	return &coll, nil
}
