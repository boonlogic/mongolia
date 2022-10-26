package mongodm

import (
	"errors"
	"fmt"
	"gitlab.boonlogic.com/development/expert/mongolia/mongodm/options"
	"go.mongodb.org/mongo-driver/mongo"
)

type ODM struct {
	mongo     *mongo.Client
	colls     map[string]*Collection
	ephemeral bool
}

func NewODM(opts *options.ODMOptions) (*ODM, error) {
	dbopts :=
	options.DBOptions().
		SetName(opts.Name).
		SetURI(opts.Host).
		db, err := mongo.Connect(ctx(), dbopts)
	if err != nil {
		return nil, err
	}
	inst := &ODM{
		db:        db,
		colls:     make(map[string]*Collection),
		ephemeral: ephemeral(opts),
	}
	return inst, nil
}

func ephemeral(opts *options.ODMOptions) bool {
	return eph
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
	schema, err := NewSchema(spec, hooks)
	if err != nil {
		return err
	}
	coll, err := NewCollection(name, schema)
	if err != nil {
		return err
	}
	o.colls[name] = coll
	return nil
}
