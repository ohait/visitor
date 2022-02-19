package service

import (
	"context"
	"encoding/json"
	"io"
	"net"
	"net/http"

	"github.com/ohait/visitor/ctx"
	"github.com/ohait/visitor/data"
)

// Will block
func (s *Service) ListenHttp(c ctx.C, addr string) error {
	mux := http.NewServeMux()
	mux.HandleFunc(`/api/add_view`, s.handleAddView)
	mux.HandleFunc(`/api/get_events`, s.handleAddView)

	l, err := net.Listen("tcp", addr)
	if err != nil {
		return ctx.Wrapf(c, "can't listen to %q: %v", addr, err)
	}
	httpd := http.Server{
		BaseContext: func(_ net.Listener) context.Context {
			return c
		},
		Handler: mux,
	}
	ctx.Log(c).Debugf("listening to %q", l.Addr().String())

	go httpd.Serve(l) // serve in the background, we can safely ignore any error

	<-ctx.Shutdown // wait for shutdown
	ctx.Log(c).Debugf("shutting down %q", l.Addr().String())

	return httpd.Shutdown(c) // blocks until all requests are finished (NOTE: hijacked might need special handling)
}

func (s *Service) handleAddView(w http.ResponseWriter, r *http.Request) {
	c := ctx.WithTag(r.Context(), `req`, r.URL.Path)
	switch r.Method {
	case http.MethodPost:
		in, err := io.ReadAll(r.Body)
		if err != nil {
			ctx.Log(c).Warnf("can't read body: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		defer r.Body.Close()
		var ev data.Event
		err = json.Unmarshal(in, &ev)
		if err != nil {
			ctx.Log(c).Warnf("invalid body: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = s.AddView(c, ev)
		if err != nil {
			ctx.Log(c).Warnf("can't add event: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

	default:
		w.WriteHeader(http.StatusBadRequest)
	}
}

func handleGetEvents(w http.ResponseWriter, r *http.Request) {
}
