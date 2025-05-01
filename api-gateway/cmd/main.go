package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/barcek2281/comics-store/api-gateway/internal/configs"
	"github.com/barcek2281/comics-store/api-gateway/internal/server"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config-path", "./configs/local.yaml", "config path")
}
func main() {
	flag.Parse()

	cfg := configs.MustLoad(configPath)

	log := slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)

	s := server.NewServer(log, cfg.Port)
	slog.Info("server starting", "port", cfg.Port)
	if err := s.Run(); err != nil {
		fmt.Printf("cannot start a server")
	}

}
