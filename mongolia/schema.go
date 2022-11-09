package mongolia

import (
	"github.com/santhosh-tekuri/jsonschema/v5"
)

type Schema struct {
	jsonschema.Schema
}

func (s *Schema) Validate(doc map[string]any) error {
	return s.Validate(doc)
}

//type Schema struct {
//	definition map[string]any
//	validator  func(map[string]any) error
//}
//
//// NewSchema creates a new Schema based on the given Spec.
//func NewSchema(spec *Spec) (*Schema, error) {
//	return newSchema(spec)
//}
//
//func newSchema(spec *Spec) (*Schema, error) {
//	vfunc, err := spec.GetValidator()
//	if err != nil {
//		return nil, err
//	}
//
//	var def map[string]any
//	if err = json.Unmarshal(spec.definition, &def); err != nil {
//		return nil, err
//	}
//
//	s := &Schema{
//		definition: def,
//		validator:  vfunc,
//	}
//	return s, nil
//}
//
//type Spec struct {
//	definition []byte
//}
//
//// NewSpec creates a new spec using the given JSON definition.
//func NewSpec(definition []byte) *Spec {
//	return newSpec(definition)
//}
//
//// NewSpecFromFile creates a new spec from the JSON definition in the given file.
//func NewSpecFromFile(path string) (*Spec, error) {
//	buf, err := ioutil.ReadFile(path)
//	if err != nil {
//		return nil, err
//	}
//	return newSpec(buf), nil
//}
//
//func newSpec(definition []byte) *Spec {
//	return &Spec{definition: definition}
//}
//
//// Validate ensures that the definition is valid.
//func (s *Spec) Validate() error {
//	return s.validate()
//}
//
//func (s *Spec) validate() error {
//	// If compiler.AddResource succeeds, the spec is valid jsonschema.
//	compiler := jsonschema.NewCompiler()
//	if err := compiler.AddResource("nil", bytes.NewBuffer(s.definition)); err != nil {
//		return err
//	}
//	return nil
//}
//
//// GetValidator returns a function that validates a map[string]any against this Spec.
//func (s *Spec) GetValidator() (func(map[string]any) error, error) {
//	return s.getValidator()
//}
//
//func (s *Spec) getValidator() (func(map[string]any) error, error) {
//	compiler := jsonschema.NewCompiler()
//	if err := compiler.AddResource("nil", bytes.NewBuffer(s.definition)); err != nil {
//		return nil, err
//	}
//	schema, err := compiler.Compile("nil")
//	if err != nil {
//		return nil, err
//	}
//	return needMap(schema.Validate), nil
//}
//
//// The signature for jsonschema.Schema's Validate accepts any type, but it panics when the JSON value is not a
//// map[string]any. Decorate the function so it requires a map[string]any.
//func needMap(fn func(any) error) func(map[string]any) error {
//	return func(v map[string]any) error {
//		if err := fn(v); err != nil {
//			return err
//		}
//		return nil
//	}
//}
