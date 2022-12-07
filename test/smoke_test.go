package test

import (
	"fmt"
	"github.com/boonlogic/mongolia"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func Test(t *testing.T) {

	odm := mongolia.NewODM().
		SetURI("mongodb://localhost:27017").
		SetDBName("mongolia-local-tmp").
		SetTimeout(10 * time.Second)

	// connect to mongo
	err := odm.Connect()
	require.Nil(t, err)

	// start clean, end clean
	odm.Drop()
	defer odm.Drop()

	// adding a schema creates a corresponding collection
	coll := odm.CreateCollection("user")
	require.NotNil(t, coll)

	// get the collection that was made
	coll = odm.GetCollection("user")
	require.NotNil(t, coll)

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

	// find the created user with FindOne
	q := map[string]any{
		"userId": uid,
	}
	err = coll.FindOne(q, found, nil)
	require.Nil(t, err)
	require.Equal(t, found.GetID(), user.GetID())

	// fail to find a user with FindOne
	found = new(User)
	q["userId"] = NewUserID()
	coll.FindOne(q, found, nil)
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

	// call Update to make the DB document match the struct
	// this updates the associated document in collection "user"
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
	// this deletes the DB document corresponding to user
	err = coll.DeleteByID(user)
	require.Nil(t, err)

	// user is gone
	found = new(User)
	err = coll.FindByID(user.GetID(), found)
	require.NotNil(t, err)
	require.Empty(t, found)
}
