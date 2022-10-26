package mongodm

import (
	"bytes"
	"github.com/santhosh-tekuri/jsonschema/v5"
	"io/ioutil"
)

type Spec struct {
	definition []byte
}

// NewSpec creates a new spec using the given JSON definition.
func NewSpec(definition []byte) *Spec {
	return newSpec(definition)
}

// NewSpecFromFile creates a new spec from the JSON definition in the given file.
func NewSpecFromFile(path string) (*Spec, error) {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return newSpec(buf), nil
}

func newSpec(definition []byte) *Spec {
	return &Spec{definition: definition}
}

// Validate ensures that the definition is valid.
func (s *Spec) Validate() error {
	return s.validate()
}

func (s *Spec) validate() error {
	// If compiler.AddResource succeeds, the spec is valid jsonschema.
	compiler := jsonschema.NewCompiler()
	if err := compiler.AddResource("nil", bytes.NewBuffer(s.definition)); err != nil {
		return err
	}
	return nil
}

// GetValidator returns the validator function for this Spec.
// The validator function determines whether an Attributes matches the definition.
func (s *Spec) GetValidator() (func(Attributes) error, error) {
	return s.getValidator()
}

func (s *Spec) getValidator() (func(Attributes) error, error) {
	compiler := jsonschema.NewCompiler()
	if err := compiler.AddResource("nil", bytes.NewBuffer(s.definition)); err != nil {
		return nil, err
	}
	schema, err := compiler.Compile("nil")
	if err != nil {
		return nil, err
	}
	vfunc := requireMapStringAny(schema.Validate)
	return decorateValidator(vfunc), nil
}

// Decorate a general-purpose validator so it accepts an Attributes.
func decorateValidator(fn func(map[string]any) error) func(Attributes) error {
	return func(v Attributes) error {
		if err := fn(v); err != nil {
			return err
		}
		return nil
	}
}

// The signature for jsonschema.Schema.Validate allows it to accept
// any type, but it panics when the JSON value is not a map[string]any.
// This decorates the function so it explicitly requires a map[string]any.
func requireMapStringAny(fn func(any) error) func(map[string]any) error {
	return func(v map[string]any) error {
		if err := fn(v); err != nil {
			return err
		}
		return nil
	}
}
