package mongodm

import (
	"context"
)

var odm *ODM

func init() {
	odm = new(ODM)
}

func ctx() context.Context {
	return context.Background()
}
