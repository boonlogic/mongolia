package test

import (
	"fmt"

	"gitlab.boonlogic.com/development/expert/mongolia/mongolia"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	mongolia.DefaultModel `bson:",inline"`
	UserID                *string              `json:"userId,omitempty" bson:"userId,omitempty"`
	Username              *string              `json:"username,omitempty" bson:"username,omitempty"`
	Perms                 []primitive.ObjectID `json:"perms,omitempty" bson:"perms,omitempty" ref:"permission"`
}

func (user *User) PreCreate() error {
	// Call the DefaultModel Creating hook
	if err := user.DefaultModel.PreCreate(); err != nil {
		return err
	}

	fmt.Println("PRECREATE HOOK")
	return nil
}

func (user *User) PreSave() error {
	// Call the DefaultModel saving hook
	if err := user.DefaultModel.PreSave(); err != nil {
		return err
	}

	fmt.Println("PRESAVE HOOK")
	return nil
}

func (user *User) GetTagReferences() map[string]string {
	return mongolia.GetStructTags(*user, "ref")
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
