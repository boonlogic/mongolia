package mongodm

import (
	"gitlab.boonlogic.com/development/expert/mongolia/mongodm/options"
	"go.mongodb.org/mongo-driver/mongo"
)

type ODM struct {
	db        *mongo.Database
	colls     map[string]*mongo.Collection
	ephemeral bool
}

func NewODM(opts *options.ConfigureOptions) (ODM, error) {
	odm.db = c.Database(*opts.Database)
	odm.colls = make(map[string]*mongo.Collection)
	odm.ephemeral = eph
}
