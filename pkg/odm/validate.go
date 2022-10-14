package odm

import (
	"errors"
	"fmt"
	"strings"

	"github.com/santhosh-tekuri/jsonschema/v5"
	"gopkg.in/yaml.v2"
)

var yamlText = `
productId: 1
productName: A green door
price: 12.50
tags:
- home
- green
`


func Validate() {
	var m interface{}
	err := yaml.Unmarshal([]byte(yamlText), &m)
	if err != nil {
		panic(err)
	}
	m, err = toStringKeys(m)
	if err != nil {
		panic(err)
	}
	compiler := jsonschema.NewCompiler()
	if err := compiler.AddResource("schema.json", strings.NewReader(schemaText)); err != nil {
		panic(err)
	}
	schema, err := compiler.Compile("schema.json")
	if err != nil {
		panic(err)
	}
	if err := schema.Validate(m); err != nil {
		panic(err)
	}
	fmt.Println("validation successfull")
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
