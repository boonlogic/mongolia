package mongodm

import (
	"errors"
	"log"
)

// Drop drops all ODM data. It will fail if the ODM is in production.
func Drop() {
	if err := drop(); err != nil {
		panic(err)
	}
}

func drop() error {
	if odm.ephemeral {
		log.Printf("mongoold '%s' dropped", odm.db.Name())
		return odm.db.Drop(ctx())
	}
	return errors.New("drop is not allowed in production mode")
}
