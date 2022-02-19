package main

import (
	"os"

	"github.com/ohait/visitor/ctx"
	"github.com/ohait/visitor/db"
	"github.com/ohait/visitor/service"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	c, cf := ctx.Root()
	defer cf()
	ctx.Log(c).Debugf("starting up...")
	defer ctx.Log(c).Debugf("exiting...")

	go ctx.WaitForSignal(cf)

	db, err := db.NewSQLStore(c, `sqlite3`, `visitor.db`) // TODO use env
	if err != nil {
		ctx.Log(c).Warnf("db: %+v", err)
		os.Exit(-1)
	}

	s := service.Service{
		DB: db,
		Auth: func(c ctx.C, token string) (service.User, error) {
			// TODO: parse token, validate, and return the right user
			return service.God{}, nil
		},
	}

	err = s.ListenHttp(c, ":8080")
	if err != nil {
		ctx.Log(c).Warnf("can't listen: %v", err)
		os.Exit(-1)
	}
}
