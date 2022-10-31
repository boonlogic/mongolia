package mongodm

import (
	"errors"
	"fmt"
	"gitlab.boonlogic.com/development/expert/mongolia/mongodm/options"
	"go.mongodb.org/mongo-driver/mongo"
	moptions "go.mongodb.org/mongo-driver/mongo/options"
)

// ODM is the object document mapper.
type ODM struct {
	name      string
	db        *mongo.Database
	colls     map[string]*Collection
	ephemeral bool
}

func NewODM(opts *options.ODMOptions) (*ODM, error) {
	client, err := connectMongo(opts)
	if err != nil {
		return nil, err
	}
	inst := &ODM{
		name:      *opts.Name,
		db:        client.Database(*opts.Name),
		colls:     make(map[string]*Collection),
		ephemeral: *opts.Ephemeral,
	}
	return inst, nil
}

func connectMongo(opts *options.ODMOptions) (*mongo.Client, error) {
	uri := opts.MongoURI()
	mopts := moptions.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx(), mopts)
	if err != nil {
		return nil, err
	}
	return client, nil
}

// AddSchema adds a schema to the ODM based on the given spec and hooks.
// Adding a Schema creates a corresponding Collection with the same name.
// The spec defines how documents are validated and the hooks allow custom
// processing to trigger at specific points in the document read/write cycle.
func AddSchema(name string, spec *Spec, hooks *Hooks) error {
	return odm.addSchema(name, spec, hooks)
}

func (o *ODM) addSchema(name string, spec *Spec, hooks *Hooks) error {
	if _, ok := o.colls[name]; ok {
		return errors.New(fmt.Sprintf("a schema named '%s' already exists", name))
	}
	if err := spec.Validate(); err != nil {
		return err
	}
	schema, err := newSchema(spec, hooks)
	if err != nil {
		return err
	}
	coll, err := newCollection(name, schema)
	if err != nil {
		return err
	}
	o.colls[name] = coll
	return nil
}
