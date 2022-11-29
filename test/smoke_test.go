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

	err := mongolia.Connect(cfg)
	require.Nil(t, err)

	mongolia.Drop()

	// no collections exist at startup
	coll, err := mongolia.GetCollection("nonexistent")
	require.NotNil(t, err)
	require.Nil(t, coll)

	// when a schema is added, a corresponding collection is made
	coll, err = mongolia.AddSchema("user", "test/user.json")
	require.Nil(t, err)

	// get the collection that was made
	coll, err = mongolia.GetCollection("user")
	require.Nil(t, err)
	require.NotNil(t, coll)

	// make a new user in memory
	var user *User
	uid := NewUserID()
	name := "luke"
	user = NewUser(uid, name)

	// create an associated record in the "tenants" collection
	err = coll.Create(user, nil)
	if err != nil {
		fmt.Printf("error: %s\n", err)
	}
	require.Nil(t, err)
	require.NotNil(t, user)
	require.Equal(t, uid, *user.UserID)
	require.Equal(t, name, *user.Username)

	// find the created user by ID
	var found = &User{}
	err = coll.FindByID(user.GetID(), found)
	require.Nil(t, err)
	require.Equal(t, found.GetID(), user.GetID())

	//// try to find a nonexistent user by ID
	//found = new(User)
	//err = coll.FindByID(NewUserID(), found)
	//require.NotNil(t, err)
	//require.
	//
	//// find the created user using First
	//q := map[string]any{
	//	"uid": uid,
	//}
	//err = coll.First(q, found, nil)
	//require.Nil(t, err)
	//require.Equal(t, found.GetID(), user.GetID())
	//
	//// try to find a nonexistent user using First
	//var result *User
	//q = map[string]any{
	//	"uid": NewUserID(),
	//}
	//coll.First(q, result, nil)
	//require.Nil(t, err)
	//require.Nil(t, result)

	// change the user in memory
	tid2 := NewUserID()
	name2 := "brad"
	user.UserID = &tid2
	user.Username = &name2

	// its DB document should remain unchanged
	err = coll.FindByID(user.GetID(), found)
	require.Nil(t, err)
	require.NotEqual(t, *found.UserID, tid2)
	require.NotEqual(t, *found.Username, name2)

	// update the DB document to match the struct
	err = coll.Update(user, nil)
	require.Nil(t, err)

	// validate that Update did not change the user struct
	require.Equal(t, *user.UserID, tid2)
	require.Equal(t, *user.Username, name2)

	// validate that the DB document matches the user struct
	err = coll.FindByID(user.GetID(), found)
	require.Nil(t, err)
	require.Equal(t, *found.UserID, *user.UserID)
	require.Equal(t, *found.Username, *user.Username)

	// delete the user
	err = coll.Delete(user)
	require.Nil(t, err)

	// make sure they're gone
	found = new(User)
	err = coll.FindByID(user.GetID(), found)
	require.NotNil(t, err)
	require.Empty(t, found)
}
