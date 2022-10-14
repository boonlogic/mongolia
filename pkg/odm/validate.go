package odm

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/santhosh-tekuri/jsonschema/v5"
)

func ValidateSpec(name string, spec []byte) (func(any) error, error) {
	url := encodeURL(name)

	compiler := jsonschema.NewCompiler()
	if err := compiler.AddResource(url, bytes.NewReader(spec)); err != nil {
		return nil, err
	}

	schema, err := compiler.Compile(url)
	if err != nil {
		return nil, err
	}

	vfunc := func(obj any) error {
		return validate(obj, schema.Validate)
	}
	return vfunc, nil
}

func validate(in any, vfunc func(any) error) error {
	buf, err := json.Marshal(in)
	if err != nil {
		return err
	}
	var obj interface{}
	if err := json.Unmarshal(buf, &obj); err != nil {
		return err
	}
	obj, err = toStringKeys(obj)
	if err != nil {
		return err
	}
	if err := vfunc(obj); err != nil {
		return err
	}
	return nil
}

func toStringKeys(val interface{}) (interface{}, error) {
	var err error
	switch val := val.(type) {
	case map[interface{}]interface{}:
		m := make(map[string]interface{})
		for k, v := range val {
			k, ok := k.(string)
			if !ok {
				return nil, errors.New(fmt.Sprintf("found non-string key: %+v", k))
			}
			m[k], err = toStringKeys(v)
			if err != nil {
				return nil, err
			}
		}
		return m, nil
	case []interface{}:
		var l = make([]interface{}, len(val))
		for i, v := range val {
			l[i], err = toStringKeys(v)
			if err != nil {
				return nil, err
			}
		}
		return l, nil
	default:
		return val, nil
	}
}
