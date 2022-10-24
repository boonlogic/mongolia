package mongodm

// Contains global singletons.

import (
	"context"
)

const (
	defaultEphemeral = true
	defaultCloud     = false
	defaultPort      = uint16(27017)
)

var (
	odm *ODM
)

func init() {
	odm = new(ODM)
}

func ctx() context.Context {
	return context.Background()
}
