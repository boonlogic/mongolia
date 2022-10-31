package mongodm

import (
	"context"
	"github.com/Kamva/mgm"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var odm *ODM

func init() {
	odm = new(ODM)
}

func ctx() context.Context {
	return context.Background()
}

func Connect(uri string, dbname string) error {
	return mgm.SetDefaultConfig(nil, dbname, options.Client().ApplyURI(uri))
}

func Coll(m *Model) *Collection {
	if c, ok := odm.colls[m.Name()]; ok {
		return c
	}
	return nil
}
