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
	inventoryHandler *handler.InventoryHandler
	orderHanler *handler.OrderHandler
}

func NewServer(log *slog.Logger, port int) *Server {
	return &Server{
		log: log,
		port: port,
		mux:  http.NewServeMux(),
		authHandler: handler.NewAuthHandler(log, 50051),
		inventoryHandler: handler.NewInventoryHandler(log, 50052),
		orderHanler: handler.NewOrderHandler(log, 50053),
	}
}

func (s *Server) Run() error {
	s.configure()

	return http.ListenAndServe(fmt.Sprintf(":%d", s.port), s.mux)
}

func (s *Server) configure() {
	s.mux.Handle("POST /auth/login", s.authHandler.Login())
	s.mux.Handle("POST /auth/register", s.authHandler.Register())

	s.mux.Handle("inventory/create", s.inventoryHandler.Create())
	s.mux.Handle("inventory/delete", s.inventoryHandler.Delete())
	s.mux.Handle("inventory/update", s.inventoryHandler.Update())
	s.mux.Handle("inventory/list", s.inventoryHandler.List())
	s.mux.Handle("inventory/get", s.inventoryHandler.Get())

	s.mux.Handle("order/create", s.orderHanler.CreateOrder())
	s.mux.Handle("order/get", s.orderHanler.GetOrder())
	s.mux.Handle("order/update", s.orderHanler.UpdateOrder())
	s.mux.Handle("order/close", s.orderHanler.CloseOrder())
	s.mux.Handle("order/list", s.orderHanler.ListOrders())
	s.mux.Handle("order/delete", s.orderHanler.DeleteOrder())
}
