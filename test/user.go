package test

import (
	"errors"
	"fmt"

	"github.com/boonlogic/mongolia/mongolia"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	mongolia.DefaultModel `bson:",inline"`
	UserID                *string              `json:"userId,omitempty" bson:"userId,omitempty"`
	Username              *string              `json:"username,omitempty" bson:"username,omitempty"`
	Perms                 []primitive.ObjectID `json:"perms,omitempty" bson:"perms,omitempty" ref:"permission"`
}

func (user *User) SetIndexes() map[string]string {
	indexes := make(map[string]string)
	indexes["username"] = "-1"
	indexes["userId"] = "1"
	return indexes
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

func (user *User) ValidateRead() error {
	fmt.Println("VALIDATE READ")
	return nil
}

func (user *User) PreUpdate(update any) error {
	err := user.DefaultModel.PreUpdate(update)
	if err != nil {
		return err
	}
	return nil
}

func (user *User) ValidateUpdate(update any) error {
	updatemap := mongolia.CastToMap(update)
	//if updating the username, make sure it meets our spec
	val, exists := updatemap["username"]
	if exists {
		username := val.(string)
		if len(username) > 64 || len(username) < 1 {
			return errors.New("Invalid username given")
		}
	}
	fmt.Println("VALIDATE UPDATE")
	return nil
}

//validation when updating entire model
func (user *User) ValidateUpdateModel() error {
	fmt.Println("VALIDATE UPDATE MODEL")
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
