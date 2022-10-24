package mongodm

import (
	"errors"
	"fmt"
	"strings"
)

// Drop drops all ODM data.
// It will fail if the database name does not end in "-tmp".
func Drop() {
	if err := drop(); err != nil {
		panic(err)
	}
}

func drop() error {
	if !strings.HasSuffix(db.Name(), "-tmp") {
		return errors.New(fmt.Sprintf("dropping database '%s' is not allowed", db.Name()))
	}
	return db.Drop(ctx())
}
