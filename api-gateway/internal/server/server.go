package server

import (
	"fmt"
	"log/slog"
	"net/http"
)

type Server struct {
	log *slog.Logger
	port   int
	mux    *http.ServeMux

}

func NewServer(log *slog.Logger, port int) *Server {
	return &Server{
		log: log,
		port: port,
		mux:  http.NewServeMux(),
	}
}

func (s *Server) Run() error {
	s.configure()

	return http.ListenAndServe(fmt.Sprintf(":%d", s.port), s.mux)
}

func (s *Server) configure() {
	s.mux.Handle("auth/login", nil)
	s.mux.Handle("auth/register", nil)

	s.mux.Handle("inventory/create", nil)
	s.mux.Handle("inventory/delete", nil)
	s.mux.Handle("inventory/update", nil)
	s.mux.Handle("inventory/list", nil)
	s.mux.Handle("inventory/get", nil)

	s.mux.Handle("order/create", nil)
	s.mux.Handle("order/get", nil)
	s.mux.Handle("order/udpate", nil)
	s.mux.Handle("order/close", nil)
	s.mux.Handle("order/list", nil)
	s.mux.Handle("order/delete", nil)
}
