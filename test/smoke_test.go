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
	err = coll.UpdateModel(user, nil)
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
	results := []*User{}
	filter := bson.M{}
	err = coll.Find(filter, &results, nil)
	require.Nil(t, err)
	require.Equal(t, len(results), 2)

	//Test FindWithResults operation
	resultswith := []*User{}
	find_results, err := coll.FindWithResults(bson.M{}, &resultswith, nil)
	require.Nil(t, err)
	require.NotNil(t, find_results)
	require.Equal(t, find_results.Filtered, int64(2))

	//Test AggregateWithResults operation
	mquery := make(map[string]string)
	matchStage := bson.D{{"$match", mquery}}
	pipeline := mongo.Pipeline{matchStage}

	projectStage := bson.D{{"$project", bson.D{{"userId", 0}}}}
	pipeline = append(pipeline, projectStage)

	sortStage := bson.D{{"$sort", bson.D{{"username", 1}}}}
	pipeline = append(pipeline, sortStage)

	skip := int64(0)
	skipStage := bson.D{{"$skip", skip}}
	pipeline = append(pipeline, skipStage)

	limit := int64(1)
	limitStage := bson.D{{"$limit", limit}}
	pipeline = append(pipeline, limitStage)

	aggregates := []*User{}
	aggregate_results, err := coll.AggregateWithResults(&aggregates, pipeline, &skip, &limit)
	require.Nil(t, err)
	require.NotNil(t, aggregate_results)
	require.Equal(t, len(aggregates), int(limit)) //test limit
	require.Equal(t, aggregate_results.Filtered, limit)
	require.Nil(t, aggregates[0].UserID)              //Test project
	require.Equal(t, *aggregates[0].Username, "brad") //Test sort

	//Test Distinct Operation
	unique_usernames, err := coll.Distinct(bson.D{}, "username")
	require.Nil(t, err)
	require.NotNil(t, unique_usernames)
	fmt.Printf("unique_usernames: %v \n", unique_usernames)

	//Test partial update
	partialusername := "stopupdating"
	partialuid := NewUserID()
	partialfilter := bson.D{{"username", name2}}
	partialupdate := make(map[string]any)
	partialupdate["username"] = partialusername
	partialupdate["userId"] = partialuid
	opts := options.Update().SetUpsert(false)
	temp := User{}
	err = coll.Update(partialfilter, partialupdate, &temp, opts)
	require.Nil(t, err)

	//Test Find operation
	refind := new(User)
	findfilter := bson.D{{"username", partialusername}}
	err = coll.FindOne(findfilter, refind, nil)
	require.Nil(t, err)
	require.Equal(t, *refind.UserID, partialuid)

	//Now test our validator (this should fail)
	//badusername
	badusername := ""
	badfilter := bson.D{{"username", partialusername}}
	badupdate := bson.D{{"$set", bson.D{
		{"username", badusername},
	},
	}}
	err = coll.Update(badfilter, badupdate, &temp, opts)
	require.NotNil(t, err)
	if err != nil {
		fmt.Printf("Expected Validate Error: %s\n", err.ToString())
	}

	//Find one and update
	fou_username := "changed_again"
	mapfilter := make(map[string]any)
	mapfilter["username"] = partialusername
	mapupdate := make(map[string]any)
	mapupdate["username"] = fou_username
	fou_model := User{}
	fou_opts := options.FindOneAndUpdate().SetUpsert(false)
	fou_opts.SetReturnDocument(options.After)
	err = coll.FindOneAndUpdate(mapfilter, mapupdate, &fou_model, fou_opts)
	require.Nil(t, err)
	require.Equal(t, *fou_model.Username, fou_username)

	// delete user
	// this deletes the DB document corresponding to user
	err = coll.Delete(user)
	require.Nil(t, err)

	// user is gone
	found = new(User)
	err = coll.FindByID(user.GetID(), found)
	if err != nil {
		fmt.Printf("Expected Error: %s\n", err.ToString())
	}
	require.NotNil(t, err)
	require.Empty(t, found)

	// delete other user
	// this deletes the DB document corresponding to name3
	deletedUser := &User{}
	query := bson.M{"username": name3}
	err = coll.DeleteOne(query, deletedUser)
	require.Nil(t, err)
	require.Equal(t, *deletedUser.Username, *user3.Username)

	//delete many
	err = coll.DeleteMany(bson.M{})
	require.Nil(t, err)

	//test struct tags
	tagUser := User{}
	refNames := mongolia.GetStructTags(tagUser, "ref")
	require.NotNil(t, refNames)
	require.Equal(t, refNames["perms"], "permission")
}
