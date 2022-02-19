package service

import (
	"github.com/ohait/visitor/ctx"
	"github.com/ohait/visitor/data"
)

type UserSession struct {
	*Service
	User User
}

func (s *UserSession) GetPerson(c ctx.C, id string) (data.Person, bool, error) {
	return s.DB.PersonById(c, id)
}

func (s *UserSession) GetViews(c ctx.C, id string) ([]data.Event, error) {
	return s.DB.EventsByPerson(c, id)
}
