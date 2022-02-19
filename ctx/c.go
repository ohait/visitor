package ctx

import (
	"context"
	"time"
)

type C context.Context // ctx.C is shorter than context.Context

func Root() (C, context.CancelFunc) {
	c := context.Background()
	return context.WithCancel(c)
}

type valC struct {
	parent C
	all    Tags
}

var _ C = valC{}

type cKey string

var all cKey = "all"

func (c valC) Value(key interface{}) interface{} {
	if key == all {
		return c.all
	}
	return c.parent.Value(key)
}

func (c valC) Deadline() (time.Time, bool) {
	return c.parent.Deadline()
}

func (c valC) Done() <-chan struct{} {
	return c.parent.Done()
}

func (c valC) Err() error {
	return c.parent.Err()
}

func getTags(c C) Tags {
	if c == nil {
		return Tags{}
	}
	v := c.Value(all)
	if v == nil {
		return Tags{}
	}
	return v.(Tags)
}

func WithTag(c C, key string, val interface{}) C {
	all := getTags(c)
	all.Set(key, val)
	return valC{
		parent: c,
		all:    all,
	}
}
