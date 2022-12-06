package mongolia

import (
	"encoding/json"

	"github.com/santhosh-tekuri/jsonschema/v5"
)

type Schema struct {
	jsonschema.Schema
	uniqueFields []string
}

func (s *Schema) Validate(model Model) error {
	buf, err := json.Marshal(model)
	if err != nil {
		return nil
	}
	var m map[string]any
	if err := json.Unmarshal(buf, &m); err != nil {
		return nil
	}
	return s.Schema.Validate(m)
}

func requiredIndexes(schema *Schema) []Index {
	idxs := make([]Index, len(schema.uniqueFields))
	for i, f := range schema.uniqueFields {
		keys := []IndexKey{
			{
				Field: f,
				Type:  Ascending,
			},
		}
		idx := Index{
			Name:   indexName(keys),
			Keys:   keys,
			Unique: true,
		}
		idxs[i] = idx
	}
	return idxs
}
