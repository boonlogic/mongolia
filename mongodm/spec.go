package mongodm

import (
	"bytes"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"io"
	"io/ioutil"
)

// Reader for the JSON definition of a schema.
type SpecReader io.Reader

func NewSpecFromFile(path string) (*SpecReader, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return loadjson(data)
}

func NewSpecFromJSON(data []byte) (*SpecReader, error) {
	return loadjson(data)
}

// Only accepts types bson.D and bson.M.
func NewSpecFromBSON(v any) (*SpecReader, error) {
	switch val := v.(type) {
	case bson.D:
		return loadbson(val)
	case bson.M:
		return loadbson(val)
	default:
		return nil, errors.New(fmt.Sprintf("value was not type bson.D or bson.M"))
	}
}

func loadbson(v any) (*SpecReader, error) {
	data, err := bson.MarshalExtJSON(v, true, false)
	if err != nil {
		return nil, err
	}
	return loadjson(data)
}

func loadjson(data []byte) (*SpecReader, error) {
	buf := bytes.NewBuffer(data)
	d := SpecReader(buf)
	return &d, nil
}
