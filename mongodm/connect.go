package mongodm

import (
	"github.com/Kamva/mgm"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var odm *ODM

func init() {
	odm = new(ODM)
}

func Connect(uri string, dbname string) error {
	return mgm.SetDefaultConfig(nil, dbname, options.Client().ApplyURI(uri))
}

func Coll(name string) *Collection {
	if c, ok := odm.colls[name]; ok {
		return c
	}
	return nil
}
