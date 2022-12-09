package test

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"
	"time"

	"github.com/boonlogic/mongolia"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Episode represents the schema for the "Episodes" collection
type Episode struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Podcast     primitive.ObjectID `bson:"podcast,omitempty"`
	Title       string             `bson:"title,omitempty"`
	Description string             `bson:"description,omitempty"`
	Duration    int32              `bson:"duration,omitempty"`
}

func Test(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://127.0.0.1:27017"))
	if err != nil {
		log.Fatal(err)
		return
	}

	collection := client.Database("quickstart").Collection("episodes")

	fileContent, err := os.Open("indextest.json")

	if err != nil {
		log.Fatal(err)
		return
	}

	fmt.Println("The File is opened successfully...")

	defer fileContent.Close()

	jsonString, _ := ioutil.ReadAll(fileContent)

	if err := mongolia.PopulateIndexes(ctx, collection, []byte(jsonString)); err != nil {
		log.Fatal(err)
		return
	}

	//all docs with asset item type
	cursor, err := collection.Find(ctx, bson.D{}, nil)
	if err != nil {
		log.Fatal(err)
		return
	}

	//Grab all matching this
	var episodes []Episode
	err = cursor.All(ctx, &episodes)
	if err != nil {
		log.Fatal(err)
		return
	}

	fmt.Printf("Found with indexes: %v \n", episodes)

}
