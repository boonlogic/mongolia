package mongodm

import (
	"encoding/json"
)

// Attributes is an unordered collection of key-value pairs.
type Attributes map[string]any

// NewAttributesFromMap creates a new Attributes using from an existing map[string]any.
func NewAttributesFromMap(v map[string]any) *Attributes {
	val := Attributes(v)
	return &val
}

// NewAttributesFromStruct creates a new Attributes from a struct using its json struct tags.
func NewAttributesFromStruct(v any) (*Attributes, error) {
	var attr Attributes
	buf, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(buf, &attr); err != nil {
		return nil, err
	}
	return &attr, nil
}
