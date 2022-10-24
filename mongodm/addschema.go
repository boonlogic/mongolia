package mongodm

import (
	"errors"
	"fmt"
	"github.com/santhosh-tekuri/jsonschema/v5"
	"gitlab.boonlogic.com/development/expert/mongolia/pkg/mongodm/options"
	"io/ioutil"
	"net/url"
)

// Singleton containing the validator functions for all schemas.
var validators = make(map[string]func(any) error)

func AddSchema(name string, spec SpecReader, hooks *Hooks) error {
	if _, ok := collections[name]; ok {
		return errors.New(fmt.Sprintf("a schema named '%s' already exists", name))
	}

	buf, err := ioutil.ReadAll(spec)
	if err != nil {
		return err
	}
	if err := registerValidator(name, buf); err != nil {
		return err
	}

	// Initialize collection here.
	// todo: add indexes
	// todo: parse indexes of x-unique attribute
	coll := options.db.Collection(name)

	s := Collection{
		collection: coll,
		validate:   validators[name],
		hooks:      hooks,
	}

	collections[name] = s
	return nil
}

func registerValidator(name string, spec SpecReader) error {
	// Compile a validator function using the jsonschema library.
	compiler := jsonschema.NewCompiler()
	url := url.QueryEscape(name)
	if err := compiler.AddResource(url, spec); err != nil {
		return err
	}
	s, err := compiler.Compile(url)
	if err != nil {
		return err
	}
	validators[name] = s.Validate
	return nil
}

type Validator func(any) error
