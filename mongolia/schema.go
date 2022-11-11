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

func (s *Schema) RequiredIndexes() []Index {
	idxs := make([]Index, len(s.uniqueFields))
	for i, f := range s.uniqueFields {
		keys := []IndexKey{
			{
				Field:     f,
				Ascending: true,
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
