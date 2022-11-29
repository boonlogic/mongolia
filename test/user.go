package test

import (
	"gitlab.boonlogic.com/development/expert/mongolia/mongolia"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	mongolia.DefaultModel `bson:",inline"`
	UserID                *string `json:"userId,omitempty" bson:"userId,omitempty"`
	Username              *string `json:"username,omitempty" bson:"username,omitempty"`
}

func NewUser(id string, name string) *User {
	return &User{
		UserID:   &id,
		Username: &name,
	}
}

func NewUserID() string {
	return primitive.NewObjectID().Hex()
}
