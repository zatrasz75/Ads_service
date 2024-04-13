package server

import (
	"net"
	"time"
)

type Option func(*Server)

func OptionSet(host, port string, redTimeout, WTimeout, iTimeout, sTimeout time.Duration) Option {
	return func(s *Server) {
		Addr(host, port)(s)
		ReadTimeout(redTimeout)(s)
		WriteTimeout(WTimeout)(s)
		IdleTimeout(iTimeout)(s)
		ShutdownTimeout(sTimeout)(s)
	}
}

func Addr(host, port string) Option {
	return func(s *Server) {
		s.server.Addr = net.JoinHostPort(host, port)
	}
}

func ReadTimeout(redTimeout time.Duration) Option {
	return func(s *Server) {
		s.server.ReadTimeout = redTimeout
	}
}

func WriteTimeout(WTimeout time.Duration) Option {
	return func(s *Server) {
		s.server.WriteTimeout = WTimeout
	}
}

func IdleTimeout(iTimeout time.Duration) Option {
	return func(s *Server) {
		s.server.IdleTimeout = iTimeout
	}
}

func ShutdownTimeout(sTimeout time.Duration) Option {
	return func(s *Server) {
		s.shutdownTimeout = sTimeout
	}
}
