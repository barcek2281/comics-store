package main

import (
	"log"
	"net"

	grpcserver "github.com/barcek2281/comics-store/auth/internal/grpcServer"
	sqlite1488 "github.com/barcek2281/comics-store/auth/internal/storage/sqlite3"
	authv1 "github.com/barcek2281/proto/gen/go/auth"
	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen on port 50051: %v", err)
	}

	store, err := sqlite1488.NewStorage("./storage/user.db")
	if err != nil {
		log.Fatalf("error to load storage: %v", err)
	}

	g := grpcserver.New(store)

	s := grpc.NewServer()
	authv1.RegisterAuthServer(s, g)

	log.Printf("gRPC server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
