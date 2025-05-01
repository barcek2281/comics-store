package main

import (
	"log/slog"
	"producer/internal/server"
)

func main() {
	
	s := server.NewServer(8181)
	slog.Info("server is running")
	if err := s.Run(); err != nil {
		slog.Error("error to running server")
	}
}