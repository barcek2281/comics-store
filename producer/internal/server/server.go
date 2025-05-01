package server

import (
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"producer/internal/mdoels"

	"github.com/nats-io/nats.go"
)

const (
	orderCreated = "order.created"
)

type Server struct {
	port int
	mux  *http.ServeMux
	nc *nats.Conn
}

func NewServer(port int) *Server {

	nc, err := nats.Connect("nats://nats:4222")
	if err != nil {
		log.Fatal(err)
	}

	return &Server{
		port: port,
		mux:  http.NewServeMux(),
		nc: nc,
	}
}

func (s *Server) Run() error {
	s.mux.HandleFunc("POST /create-order", s.createOrder())
	slog.Info(fmt.Sprintf(":%d", s.port))
	return http.ListenAndServe(fmt.Sprintf(":%d", s.port), s.mux)
}

func (s *Server) createOrder() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var orderReq mdoels.Order
		if err := json.NewDecoder(r.Body).Decode(&orderReq); err != nil {
			slog.Error("cannot parse json file")
			return
		}
		data, _ := json.Marshal(orderReq)
		err := s.nc.Publish(orderCreated, data)

		if err != nil {
			slog.Error("cannot publish", "error", err)
		}
		
		slog.Info("publish order.created")
	}
}
