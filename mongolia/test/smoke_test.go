package test

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"gitlab.boonlogic.com/development/expert/mongolia/mongolia"
	"testing"
	"time"
)

type Book struct {
	mongolia.DefaultModel `bson:",inline"`
	Title                 *string `json:"title,omitempty"`
	Author                *string `json:"name,omitempty"`
	Year                  *int    `json:"year,omitempty"`
}

func NewBook(title string, author string, year int) *Book {
	return &Book{
		Title:  &title,
		Author: &author,
		Year:   &year,
	}
}

func Test(t *testing.T) {
	cfg := mongolia.NewConfig().
		SetURI("mongodb://localhost:27017").
		SetDBName("mongolia-local").
		SetTimeout(10 * time.Second)

	err := mongolia.Connect(cfg)
	require.Nil(t, err)

	err = mongolia.AddSchema("book", "test/book.json")
	require.Nil(t, err)

	coll, err := mongolia.GetCollection("tenant")
	require.Nil(t, err)
	require.NotNil(t, coll)

	title := "The book of Amber"
	author := "J.R.R. Turnquist"
	year := 2020
	book := NewBook(title, author, year)

	// create a book
	err = coll.Create(book, nil)
	if err != nil {
		fmt.Printf("error: %s\n", err)
	}
	require.Nil(t, err)
	require.NotNil(t, book)
	require.Equal(t, title, book.Title)
	require.Equal(t, author, book.Author)
	require.Equal(t, year, book.Year)

	// change the book in memory, but not the database
	newtitle := "The book of AVIS"
	newauthor := "I.A. & B.P.T. (pseudononymous)"
	newyear := 2017
	book.Title = &newtitle
	book.Author = &newauthor
	book.Year = &newyear

	// query the created book
	var found *Book
	query := map[string]any{
		"title":  "The book of Amber",
		"author": "J.R.R. Turnquist",
		"year":   1215,
	}
	err = coll.First(query, found, nil) // unsuccessful
	require.NotNil(t, err)
	require.Nil(t, found)

	query["year"] = 2020
	err = coll.First(query, found, nil) // successful
	require.Nil(t, err)
	require.Nil(t, found)
}
