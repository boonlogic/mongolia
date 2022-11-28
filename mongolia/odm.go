package mongolia

import (
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"io/ioutil"
)

type ODM struct {
	client   *mongo.Client
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
	odm.config = config
	return nil
}

func (o *ODM) AddSchema(name string, path string) (*Collection, error) {
	def, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	schema, err := odm.compileSchema(path, def)
	if err != nil {
		return nil, err
	}
	coll, err := connectCollection(name, schema)
	if err != nil {
		return nil, err
	}
	odm.colls[name] = *coll
	return coll, nil
}

func (o *ODM) GetCollection(name string) (*Collection, error) {
	coll, ok := o.colls[name]
	if !ok {
		return nil, errors.New(fmt.Sprintf("no collection named \"%s\"", name))
	}
	return &coll, nil
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

func (o *ODM) drop() {
	_ = o.config
	if !*o.config.Ephemeral {
		panic("odm instance is not ephemeral")
	}
	for _, coll := range o.colls {
		if err := coll.drop(); err != nil {
			panic(err)
		}
	}
	// todo: drop odm database itself
}
