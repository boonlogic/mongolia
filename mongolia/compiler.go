package mongolia

import (
	"bytes"
	"github.com/santhosh-tekuri/jsonschema/v5"
)

type Compiler struct {
	jsonschema.Compiler
}

func NewCompiler() *Compiler {
	return &Compiler{
		Compiler: *jsonschema.NewCompiler(),
	}
}

func (c *Compiler) AddResource(uri string, definition []byte) error {
	c.Compiler.AddResource(uri, bytes.NewReader(definition))
	return nil
}

func (c *Compiler) Compile(uri string) (*Schema, error) {
	schema, err := c.Compiler.Compile(uri)
	if err != nil {
		return nil, err
	}
	s := &Schema{
		Schema: *schema,
	}
	return s, nil
}
