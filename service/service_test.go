package service

import (
	"testing"

	"github.com/ohait/visitor/ctx"
	"github.com/ohait/visitor/ctx/test"
	"github.com/ohait/visitor/data"
	"github.com/ohait/visitor/db"
	"github.com/stretchr/testify/require"

	_ "github.com/mattn/go-sqlite3"
)

func TestAll(t *testing.T) {
	test.Run(t, func(c test.C) {
		db, err := db.NewSQLStore(c, "sqlite3", t.TempDir()+"/"+t.Name()+".sqlite")
		c.NoError(err)
		require.NotEmpty(t, db)
		s := Service{
			Auth: func(c ctx.C, token string) (User, error) { return God{}, nil },
			DB:   db,
		}

		t.Logf("add click")
		err = s.AddView(c, data.Event{
			Campaign: "c1",
			Person:   "p1",
			Type:     data.EventTypeLandingPage,
		})
		c.NoError(err)

		t.Logf("add another click")
		err = s.AddView(c, data.Event{
			Campaign: "c1",
			Person:   "p1",
			Type:     data.EventTypeVideoPlay,
		})
		c.NoError(err)

		t.Logf("login")
		{
			s, err := s.NewSession(c, "god")
			c.NoError(err)

			events, err := s.GetViews(c, "p1")
			c.NoError(err)
			t.Logf("got events: %+v", events)
			c.Equal(2, len(events))
		}
	})
}
