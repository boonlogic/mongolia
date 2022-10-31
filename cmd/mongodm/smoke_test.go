package main

import (
	"gitlab.boonlogic.com/development/expert/mongolia/mongodm"
	"testing"
)

type Album struct {
	mongodm.DefaultModel `bson:",inline"`
	Name                 string `json:"name" bson:"name"`
	Artist               string `json:"artist" bson:"artist"`
}

func NewAlbum(name string, artist string) *Album {
	return &Album{
		Name:   name,
		Artist: artist,
	}
}

func TestSmoke(t *testing.T) {
	mongodm.Connect("mongodb://localhost:27017", "mongodm")

	album := NewAlbum("kind of blue", "miles davis")
	Coll()
}
