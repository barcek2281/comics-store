package main

import (
	"log"
	"net"

	grpcserver "github.com/barcek2281/comics-store/inventory/internal/grpcServer"
	"github.com/barcek2281/comics-store/inventory/internal/storage/sqlite"
	inventoryv1 "github.com/barcek2281/proto/gen/go/inventory"
	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("failed to listen on port 50052: %v", err)
	}

	store, err := sqlite.NewStorage("./storage/user.db")
	if err != nil {
		log.Fatalf("error to load storage: %v", err)
	}

	g := grpcserver.New(store)

	s := grpc.NewServer()
	inventoryv1.RegisterInventoryServer(s, g)

	log.Printf("gRPC server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
