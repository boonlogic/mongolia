package odm

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/santhosh-tekuri/jsonschema/v5"
)

func ValidateJSONSchema(value []byte, spec []byte) (bool, error) {
	var obj interface{}
	if err := json.Unmarshal(value, &obj); err != nil {
		return false, err
	}

	obj, err := toStringKeys(obj)
	if err != nil {
		return false, err
	}

	compiler := jsonschema.NewCompiler()
	if err := compiler.AddResource("schema.json", bytes.NewReader(spec)); err != nil {
		return false, err
	}

	schema, err := compiler.Compile("schema.json")
	if err != nil {
		return false, err
	}

	if err := schema.Validate(value); err != nil {
		return false, err
	}
	return true, nil
}

func toStringKeys(val interface{}) (interface{}, error) {
	var err error
	switch val := val.(type) {
	case map[interface{}]interface{}:
		m := make(map[string]interface{})
		for k, v := range val {
			k, ok := k.(string)
			if !ok {
				return nil, errors.New("found non-string key")
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
