package test

import (
	"context"
	"fmt"

	"gitlab.boonlogic.com/development/expert/mongolia/mongolia"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	mongolia.DefaultModel `bson:",inline"`
	UserID                *string `json:"userId,omitempty" bson:"userId,omitempty"`
	Username              *string `json:"username,omitempty" bson:"username,omitempty"`
}

func (user *User) PreCreate(ctx context.Context) error {
	// Call the DefaultModel Creating hook
	if err := user.DefaultModel.PreCreate(ctx); err != nil {
		return err
	}

	fmt.Println("PRECREATE HOOK")
	return nil
}

func (user *User) PreSave(ctx context.Context) error {
	// Call the DefaultModel saving hook
	if err := user.DefaultModel.PreSave(ctx); err != nil {
		return err
	}

	fmt.Println("PRESAVE HOOK")
	return nil
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
