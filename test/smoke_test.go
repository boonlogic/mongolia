package test

import (
	"fmt"
	"github.com/boonlogic/mongolia/mongolia"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
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
	coll := odm.CreateCollection("user", nil)
	require.NotNil(t, coll)

	//Create Indexes
	indexes := []mongo.IndexModel{
		{
			Keys: bsonx.Doc{
				{Key: "userId", Value: bsonx.String("text")},
				{Key: "username", Value: bsonx.String("text")},
			},
			Options: options.Index().SetName("id_name").SetUnique(false),
		},
	}
	err = coll.CreateIndexes(indexes)
	require.Nil(t, err)

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
		fmt.Printf("error: %s\n", err.ToString())
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

	fmt.Printf("FindOne: %v \n", *found)

	// fail to find a user with FindOne
	found = new(User)
	q["userId"] = NewUserID()
	coll.FindOne(q, found, nil)
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

	//Add another user
	uid3 := NewUserID()
	name3 := "frank"
	user3 := NewUser(uid3, name3)

	// create the user in ODM
	// this is where the document is actually added to collection "user")
	err = coll.Create(user3, nil)
	if err != nil {
		fmt.Printf("error: %s\n", err.ToString())
	}
	require.Nil(t, err)
	require.NotNil(t, user3)
	require.Equal(t, uid3, *user3.UserID)
	require.Equal(t, name3, *user3.Username)

	//Test Find operation
	results := []User{}
	filter := bson.M{}
	find_results, err := coll.Find(filter, &results, nil)
	require.Nil(t, err)
	require.NotNil(t, find_results)
	require.Equal(t, len(results), 2)

	fmt.Printf("Find: %v \n", results)

	//Test Distinct Operation
	unique_usernames, err := coll.Distinct(bson.D{}, "username")
	require.Nil(t, err)
	require.NotNil(t, unique_usernames)
	fmt.Printf("unique_usernames: %v \n", unique_usernames)

	// delete user
	// this deletes the DB document corresponding to user
	err = coll.DeleteModel(user)
	require.Nil(t, err)

	// user is gone
	found = new(User)
	err = coll.FindByID(user.GetID(), found)
	if err != nil {
		fmt.Printf("Expected Error: %s\n", err.ToString())
	}
	require.NotNil(t, err)
	require.Empty(t, found)
}
