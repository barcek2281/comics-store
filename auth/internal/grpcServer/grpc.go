package grpcserver

import (
	"context"

	 "github.com/barcek2281/comics-store/auth/internal/storage/sqlite3"
	authv1 "github.com/barcek2281/proto/gen/go/auth"
)

type GRPCserver struct {
	store *sqlite1488.Storage
	authv1.UnimplementedAuthServer
}

func New(store *sqlite1488.Storage) *GRPCserver {
	return &GRPCserver{
		store: store,
	}
}

func (g *GRPCserver) Register(ctx context.Context, in *authv1.RegisterRequest) (*authv1.RegisterResponse, error) {
	return &authv1.RegisterResponse{Token: "lox"}, nil
}

func (g *GRPCserver) Login(ctx context.Context, in *authv1.LoginRequest) (*authv1.LoginResponse, error) {
	return &authv1.LoginResponse{Token: "lox"}, nil
}
