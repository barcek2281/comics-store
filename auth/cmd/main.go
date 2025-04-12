package main

import (
	"log"
	"net"

	grpcserver "github.com/barcek2281/comics-store/auth/internal/grpcServer"
	authv1 "github.com/barcek2281/proto/gen/go/auth"
	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen on port 50051: %v", err)
	}
	s := grpc.NewServer()
	authv1.RegisterAuthServer(s, &grpcserver.GRPCserver{})

	log.Printf("gRPC server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
