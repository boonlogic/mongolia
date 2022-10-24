package mongodm

import "go.mongodb.org/mongo-driver/mongo"

type FindOptions struct {
	Joins      *Joins
	Projection *Projection
	Skip       *Skip
	Limit      *Limit
	Sort       *Sort
}

type Joins map[string]bool
type Projection map[string]bool
type Skip bool
type Limit int
type Sort struct {
	dir bool
}

type RemoveOptions struct {
	Projection *Projection
	Skip       *Skip
	Limit      *Limit
	Sort       *Sort
}

type Cursor struct {
	cur *mongo.Cursor
}

type Query interface {
}

type FindQuery struct{}
type RemoveQuery struct{}

func Find(query Query)
func Remove(query Query)

func find(query FindQuery) (Cursor, error) {
}

func remove(query RemoveQuery) (Cursor, error) {
}
