package mongodm

import (
	"go.mongodb.org/mongo-driver/bson"
)

type Document struct {
	attrs Attributes
}

func (d Document) BSON() []byte {
	buf, err := bson.Marshal(d.attrs)
	if err != nil {
		return nil
	}
	return buf
}
