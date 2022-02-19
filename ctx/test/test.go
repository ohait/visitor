package test

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/ohait/visitor/ctx"
)

type C struct {
	ctx.C
	t *testing.T
}

func Run(t *testing.T, fn func(c C)) {
	c := context.Background()
	c, cf := context.WithTimeout(c, 5*time.Second) // we leak the
	defer cf()
	fn(C{c, t})
}

func jsony(obj interface{}) string {
	switch obj := obj.(type) {
	case json.RawMessage:
		return string(obj)
	case string:
		return obj
	}
	j, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		return err.Error()
	}
	return string(j)
}

func (c C) NoError(err error) {
	c.t.Helper()
	if err != nil {
		var e ctx.Error
		if errors.As(err, &e) {
			c.t.Fatalf(jsony(map[string]interface{}{
				"err":   err.Error(),
				"stack": e.Stack,
				"tags":  e.Tags,
			}))
		}
		c.t.Fatalf("%s", err)
	}
}

func (c C) Equal(expected, got interface{}) {
	e := jsony(expected)
	g := jsony(got)
	if e != g {
		c.t.Fatalf("unexpected: %s", g)
	}
}
