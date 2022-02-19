package service

import (
	"github.com/ohait/visitor/ctx"
	"github.com/ohait/visitor/data"
	"github.com/ohait/visitor/db"
)

type Service struct {
	DB   db.Store
	Auth Auth
}

func (s *Service) NewSession(c ctx.C, token string) (UserSession, error) {
	u, err := s.Auth(c, token)
	return UserSession{s, u}, err
}

func (s *Service) AddView(c ctx.C, ev data.Event) error {
	if ev.Url != "" {
		// TODO expand URL into person id and campaing id
	}
	if ev.Campaign == "" {
		return ctx.Wrapf(c, "campaing missing")
	}
	if ev.Person == "" {
		return ctx.Wrapf(c, "person missing")
	}
	return s.DB.InsertEvent(c, ev)
}
