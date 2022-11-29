package test

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"gitlab.boonlogic.com/development/expert/mongolia/mongolia"
	"testing"
	"time"
)

func Test(t *testing.T) {
	cfg := mongolia.NewConfig().
		SetURI("mongodb://localhost:27017").
		SetDBName("mongolia-local-tmp").
		SetTimeout(10 * time.Second).
		SetEphemeral(true)

	// connect to mongo
	err := mongolia.Connect(cfg)
	require.Nil(t, err)

	// start clean, end clean
	mongolia.Drop()
	defer mongolia.Drop()

	// adding a schema creates a corresponding collection
	coll, err := mongolia.AddSchema("user", "test/user.json")
	require.Nil(t, err)

	// get the collection that was made
	coll, err = mongolia.GetCollection("user")
	require.Nil(t, err)
	require.NotNil(t, coll)

	// fail to get a collection
	result, err := mongolia.GetCollection("nonexistent")
	require.NotNil(t, err)
	require.Nil(t, result)

	// make a new user struct in memory
	uid := NewUserID()
	name := "luke"
	user := NewUser(uid, name)

	// create the user in ODM
	// this is where the document is actually added to collection "user")
	err = coll.Create(user, nil)
	if err != nil {
		fmt.Printf("error: %s\n", err)
	}
	require.Nil(t, err)
	require.NotNil(t, user)
	require.Equal(t, uid, *user.UserID)
	require.Equal(t, name, *user.Username)

	// find the created user by ID
	var found = new(User)
	err = coll.FindByID(user.GetID(), found)
	require.Nil(t, err)
	require.Equal(t, found.GetID(), user.GetID())

	// fail to find a user by ID
	found = new(User)
	err = coll.FindByID(NewUserID(), found)
	require.NotNil(t, err)
	require.Empty(t, found)

	// find the created user with First
	q := map[string]any{
		"userId": uid,
	}
	err = coll.First(q, found, nil)
	require.Nil(t, err)
	require.Equal(t, found.GetID(), user.GetID())

	// fail to find a user with First
	found = new(User)
	q["userId"] = NewUserID()
	coll.First(q, found, nil)
	require.Nil(t, err)
	require.Empty(t, found)

	// change the user struct in memory
	tid2 := NewUserID()
	name2 := "brad"
	user.UserID = &tid2
	user.Username = &name2

	// note that the DB document remains unchanged
	err = coll.FindByID(user.GetID(), found)
	require.Nil(t, err)
	require.NotEqual(t, *found.UserID, tid2)
	require.NotEqual(t, *found.Username, name2)

	// update the DB document to match the struct
	err = coll.Update(user, nil)
	require.Nil(t, err)

	// calling Update does not change the struct in memory
	require.Equal(t, *user.UserID, tid2)
	require.Equal(t, *user.Username, name2)

	// the DB document now matches the user struct
	err = coll.FindByID(user.GetID(), found)
	require.Nil(t, err)
	require.Equal(t, *found.UserID, *user.UserID)
	require.Equal(t, *found.Username, *user.Username)

	// delete user
	err = coll.Delete(user)
	require.Nil(t, err)

	// make sure they're gone
	found = new(User)
	err = coll.FindByID(user.GetID(), found)
	require.NotNil(t, err)
	require.Empty(t, found)
}
