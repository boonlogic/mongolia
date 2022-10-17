package odm

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/santhosh-tekuri/jsonschema/v5"
	"net/url"
)

// todo: factor out compiler into a singleton
// todo: separate func that after validation/registry is done, returns the decorated validator funcs (it has to get those off the singleton)

// Ensure that the jsonschema spec for a odm type is valid.
// Produces a specialized validator func for this odm type.
func validateSpec(name string, spec []byte) (func(any) error, error) {
	url := url.QueryEscape(name)
	compiler := jsonschema.NewCompiler()
	if err := compiler.AddResource(url, bytes.NewReader(spec)); err != nil {
		return nil, err
	}
	schema, err := compiler.Compile(url)
	if err != nil {
		return nil, err
	}
	return makeValidator(schema.Validate), nil
}

// Decorate validator function with signature that calls it against `obj`.
func makeValidator(validator func(any) error) (func(any) error) {
	return func(obj any) error {
		return validateStructWithFunc(obj, validator)
	}
}

// Call validator function against struct
func validateStructWithFunc(v any, validator func(any) error) error {
	obj, err := structToInterface(v)
	if err != nil {
		return err
	}
	if obj, err = keysToString(obj); err != nil {
		return err
	}
	if err := validator(obj); err != nil {
		return err
	}
	return nil
}

func structToInterface(in any) (any, error) {
	s, err := json.Marshal(in)
	if err != nil {
		return nil, err
	}
	var out any
	if err := json.Unmarshal(s, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func keysToString(in any) (any, error) {
	var err error
	switch v := in.(type) {
	case map[any]any:
		m := make(map[string]any)
		for k, v := range v {
			k, ok := k.(string)
			if !ok {
				return nil, errors.New(fmt.Sprintf("found non-string key: %+v", k))
			}
			m[k], err = keysToString(v)
			if err != nil {
				return nil, err
			}
		}
		return m, nil
	case []any:
		l := make([]any, len(v))
		for i, v := range v {
			l[i], err = keysToString(v)
			if err != nil {
				return nil, err
			}
		}
		return l, nil
	default:
		return v, nil
	}
}
