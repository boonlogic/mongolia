package types

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
)

// Attributes is an unordered collection of key-value pairs.
type Attributes map[string]any

func NewAttributesFromMap(v map[string]any) *Attributes {
	val := Attributes(v)
	return &val
}

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

// The order of fields in the resulting bson.D is undefined.
func (a Attributes) D() bson.D {
	d := make(bson.D, 0)
	for k, v := range a.M() {
		e := bson.E{
			Key:   k,
			Value: v,
		}
		d = append(d, e)
	}
	return d
}

func (a Attributes) M() bson.M {
	return bson.M(a)
}
