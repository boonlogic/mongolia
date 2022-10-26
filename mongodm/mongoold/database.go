package mongoold

import (
	driver "go.mongodb.org/mongo-driver/mongo"
)

type DB struct {
	client *driver.Client
	colls  map[string]*Collection
}

func (d *DB) Collection(name string) *Collection {
	v, ok := d.colls[name]
	if !ok {
		return nil
	}
	return v
}
