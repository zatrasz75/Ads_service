package server

import (
	"context"
	"net/http"
	"time"
)

type Server struct {
	server          *http.Server
	notify          chan error
	shutdownTimeout time.Duration
}

func New(router http.Handler, opts ...Option) *Server {
	httpserver := &http.Server{
		Handler: router,
	}

	s := &Server{
		server: httpserver,
		notify: make(chan error, 1),
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

func (s *Server) Start() error {
	err := s.server.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}

func (s *Server) Notify() <-chan error {
	return s.notify
}

func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()

	return s.server.Shutdown(ctx)
}
