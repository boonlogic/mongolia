package mongodm

import "go.mongodb.org/mongo-driver/bson/primitive"

type Document struct {
	id    primitive.ObjectID
	attrs *Attributes
}
