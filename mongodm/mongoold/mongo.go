package mongoold

import "context"

var db *DB

func init() {
	db = new(DB)
}

func ctx() context.Context {
	return context.Background()
}
