package main

import (
	"gitlab.boonlogic.com/development/expert/mongolia/mongolia"
	"testing"
)

type Album struct {
	mongolia.DefaultModel `bson:",inline"`
	Name                  string `json:"name" bson:"name"`
	Artist                string `json:"artist" bson:"artist"`
}

func NewAlbum(name string, artist string) *Album {
	return &Album{
		Name:   name,
		Artist: artist,
	}
}

func TestSmoke(t *testing.T) {
	mongolia.Connect("mongodb://localhost:27017", "mongolia")

	album := NewAlbum("kind of blue", "miles davis")
	Coll()
}
