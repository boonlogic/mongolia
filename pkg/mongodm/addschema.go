package mongodm

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/santhosh-tekuri/jsonschema/v5"
	"go.mongodb.org/mongo-driver/bson"
	"net/url"
)

// Singleton containing the validator functions for all schemas.
var validators = make(map[string]func(any) error)

func AddSchema(name string, spec []byte, hooks *Hooks) error {
	if _, ok := collections[name]; ok {
		return errors.New(fmt.Sprintf("a schema named '%s' already exists", name))
	}

	if err := registerValidator(name, spec); err != nil {
		return nil
	}

	// Initialize collection here.
	// todo: add indexes
	// todo: parse indexes of x-unique attribute
	coll := db.Collection(name)

	s := Schema{
		collection: coll,
		validate:   validators[name],
		hooks:      hooks,
	}

	collections[name] = s
	return nil
}

func registerValidator(name string, spec []byte) error {
	// Compile a validator function using the jsonschema library.
	compiler := jsonschema.NewCompiler()
	url := url.QueryEscape(name)
	if err := compiler.AddResource(url, bytes.NewReader(spec)); err != nil {
		return err
	}
	s, err := compiler.Compile(url)
	if err != nil {
		return err
	}

	// Decorate the jsonschema validator so it accepts *Document instead of JSON data.
	fn := func(val any) error {
		doc, ok := val.(*Document)
		if !ok {
			return errors.New(fmt.Sprintf("cannot validate type %T", val))
		}
		buf, err := bson.Marshal(doc)
		if err != nil {
			return err
		}
		return s.Validate(buf)
	}

	validators[name] = fn
	return nil
}
