package mgm

import "testing"

type Album struct {
	DefaultModel `bson:",inline"`
	Name         string `json:"name" bson:"name"`
	Artist       string `json:"artist" bson:"artist"`
}

func NewAlbum(name string, artist string) *Album {
	return &Album{
		Name:   name,
		Artist: artist,
	}
}

func Test(t *testing.T) {
	album := NewAlbum("kind of blue", "miles davis")
	albumColl := Collection
}
