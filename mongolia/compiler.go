package mongolia

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/santhosh-tekuri/jsonschema/v5"
)

type Compiler struct {
	jsonschema.Compiler
	resources map[string]*Resource
}

type Resource struct {
	URI        string
	Definition Definition
}

func (r *Resource) Validate() error {
	return r.validate()
}

func (r *Resource) validate() error {
	// At the moment, the only constraint on our specification is that it is valid jsonschema.
	// We validate this by checking if it can be compiled by the jsonschema compiler.
	tester := jsonschema.NewCompiler()
	err := tester.AddResource(r.URI, bytes.NewReader(r.Definition))
	if err != nil {
		return err
	}

	// todo
	// As our specification grows, other validation constraints must be added here.
	// For example, we will need to validate the syntax of our custom x-attributes.

	return nil
}

func NewCompiler() *Compiler {
	return &Compiler{
		Compiler:  *jsonschema.NewCompiler(),
		resources: make(map[string]*Resource),
	}
}

func (c *Compiler) AddResource(name string, resource *Resource) error {
	if err := resource.validate(); err != nil {
		return err
	}
	return addResource(c, name, resource)
}

func addResource(c *Compiler, name string, r *Resource) error {
	err := c.Compiler.AddResource(r.URI, bytes.NewReader(r.Definition))
	if err != nil {
		return err
	}
	c.resources[name] = r
	return nil
}

func (c *Compiler) Compile(name string) (*Schema, error) {
	r, ok := c.resources[name]
	if !ok {
		return nil, errors.New(fmt.Sprintf("no compiler resource named \"%s\"", name))
	}
	schema, err := compileSchema(c, r)
	if err != nil {
		return nil, err
	}
	return schema, nil
}

func compileSchema(c *Compiler, r *Resource) (*Schema, error) {
	schema, err := c.Compiler.Compile(r.URI)
	if err != nil {
		return nil, err
	}
	s := &Schema{
		Schema:       *schema,
		uniqueFields: r.Definition.UniqueFields(),
	}
	return s, nil
}
