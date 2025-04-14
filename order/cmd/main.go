package main

import (
	"log"
	"net"

	"github.com/barcek2281/comics-store/order/internal/server"
	"github.com/barcek2281/comics-store/order/internal/storage"
	orderv1 "github.com/barcek2281/proto/gen/go/order"
	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":50053")
	if err != nil {
		log.Fatalf("failed to listen on port 50053: %v", err)
	}

	store, err := storage.NewStorage("../storage/database.db")
	if err != nil {
		log.Fatalf("error to load storage: %v", err)
	}

	g := server.NewGRPCserver(store)

	s := grpc.NewServer()
	orderv1.RegisterOrderServer(s, g)

	log.Printf("gRPC server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
