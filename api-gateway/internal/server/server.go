package server

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/barcek2281/comics-store/api-gateway/internal/handler"
)

type Server struct {
	log *slog.Logger
	port   int
	mux    *http.ServeMux
	authHandler *handler.AuthHandler
}

func NewServer(log *slog.Logger, port int) *Server {
	return &Server{
		log: log,
		port: port,
		mux:  http.NewServeMux(),
		authHandler: handler.NewAuthHandler(log, 50051),
	}
}

func (s *Server) Run() error {
	s.configure()

	return http.ListenAndServe(fmt.Sprintf(":%d", s.port), s.mux)
}

func (s *Server) configure() {
	s.mux.Handle("POST /auth/login", s.authHandler.Login())
	s.mux.Handle("POST /auth/register", s.authHandler.Register())

	// s.mux.Handle("inventory/create", nil)
	// s.mux.Handle("inventory/delete", nil)
	// s.mux.Handle("inventory/update", nil)
	// s.mux.Handle("inventory/list", nil)
	// s.mux.Handle("inventory/get", nil)

	// s.mux.Handle("order/create", nil)
	// s.mux.Handle("order/get", nil)
	// s.mux.Handle("order/udpate", nil)
	// s.mux.Handle("order/close", nil)
	// s.mux.Handle("order/list", nil)
	// s.mux.Handle("order/delete", nil)
}
