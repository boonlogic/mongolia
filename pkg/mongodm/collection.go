package mongodm

type Collection struct {
	name   string
	schema *Schema
}

func newCollection(name string, schema *Schema) *Collection {
	return &Collection{
		name:   name,
		schema: schema,
	}
}

// Name returns the name of the collection.
func (c *Collection) Name() string {
	return c.name
}
