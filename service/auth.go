package service

import "github.com/ohait/visitor/ctx"

type Auth func(c ctx.C, token string) (User, error)

type User interface {
	CanAccessPerson(id string) error
	CanAccessCampaign(id string) error
}

type God struct{}

func (god God) CanAccessPerson(_ string) error {
	return nil
}

func (god God) CanAccessCampaign(_ string) error {
	return nil
}
