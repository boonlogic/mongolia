package mongolia

import (
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"io/ioutil"
	"log"
)

type ODM struct {
	client   *mongo.Client
	db       *mongo.Database
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

// Connect establishes ODM's connection to mongo.
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
	log.Printf("connected to %s: '%s'", *config.URI, *config.DBName)
	return nil
}

// AddSchema adds a new Schema to ODM.
// Adding a Schema creates a corresponding Collection with the same name.
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
	log.Printf("added collection '%s'", name)
	return coll, nil
}

// GetCollection returns a Collection by name.
func (o *ODM) GetCollection(name string) (*Collection, error) {
	coll, ok := o.colls[name]
	if !ok {
		return nil, errors.New(fmt.Sprintf("no collection named \"%s\"", name))
	}
	return &coll, nil
}

// Drop deletes the ODM database.
func (o *ODM) Drop() {
	if err := o.drop(); err != nil {
		panic(fmt.Sprintf("drop failed: %s", err))
	}
	log.Printf("instance '%s' %s", *o.config.DBName, dropped())
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

func (o *ODM) drop() error {
	_ = o.config
	if !*o.config.Ephemeral {
		return errors.New("instance is not ephemeral")
	}
	return o.db.Drop(ctx())
}
