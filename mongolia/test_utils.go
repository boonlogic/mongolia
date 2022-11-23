package mongolia

import "go.mongodb.org/mongo-driver/bson/primitive"

func drop() {
	odm.drop()
}

func newoid() string {
	return primitive.NewObjectID().Hex()
}
