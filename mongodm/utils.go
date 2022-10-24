package mongodm

import (
	"errors"
)

// Drop drops all ODM data. It will fail if the ODM is not ephemeral.
func Drop() {
	if err := drop(); err != nil {
		panic(err)
	}
}

func drop() error {
	if odm.ephemeral {
		return odm.db.Drop(ctx())
	}
	return errors.New("not ephemeral, drop is not allowed")
}
