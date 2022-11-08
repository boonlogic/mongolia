package mongolia

import (
	"context"
)

var odm *ODM

func init() {
	odm = NewODM()
}

func ctx() context.Context {
	return context.Background()
}
