package mongodm

import "context"

var odm *ODM

func ctx() context.Context {
	return context.Background()
}
