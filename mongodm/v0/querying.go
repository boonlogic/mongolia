package v0

type Query struct {
}

type Cursor struct {
}

type FindOptions struct {
	Joins      *map[string]bool
	Projection *map[string]bool
	Skip       *uint64
	Limit      *uint64
	Sort       *bool
}

func find(query *Query) (*Cursor, error) {
	return nil, nil
}
