package db

import (
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/ohait/visitor/ctx/test"
	"github.com/ohait/visitor/data"
)

func TestDB(t *testing.T) {
	test.Run(t, func(c test.C) {
		dbfile := "/tmp/" + t.Name() + ".sqlite"
		os.Remove(dbfile)
		db, err := NewSQLStore(c, "sqlite3", dbfile)
		c.NoError(err)
		_, found, err := db.PersonById(c, "123")
		c.NoError(err)
		c.Equal(false, found)

		err = db.UpsertPerson(c, data.Person{
			Id:    "123",
			Email: "a@b.c",
			Details: data.Details{
				Name: "Alice",
				Addresses: []string{
					"abc road, somewhere, far away",
				},
				Phone: []string{
					"+47 123 12 123",
				},
			},
		})
		c.NoError(err)

		p, found, err := db.PersonById(c, "123")
		c.NoError(err)
		c.Equal(true, found)
		c.Equal("123", p.Id)
		c.Equal("a@b.c", p.Email)
		c.Equal("Alice", p.Details.Name)

		err = db.InsertEvent(c, data.Event{
			Person:   "123",
			Campaign: "c1",
			Type:     data.EventTypeLandingPage,
		})
		c.NoError(err)

		list, err := db.EventsByPerson(c, "123")
		c.NoError(err)
		t.Logf("events: %+v", list)
		c.Equal(1, len(list))
	})
}
