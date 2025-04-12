package grpcserver

import (
	"context"
	"time"

	"github.com/barcek2281/comics-store/auth/internal/lib/jwt"
	"github.com/barcek2281/comics-store/auth/internal/model"
	sqlite1488 "github.com/barcek2281/comics-store/auth/internal/storage/sqlite3"
	authv1 "github.com/barcek2281/proto/gen/go/auth"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

	HashPassword, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)

	user := model.User{
		Email:    in.Email,
		Password: string(HashPassword),
	}
	id, err := g.store.Save(user)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "email is used")
	}
	user.ID = id

	token, err := jwt.NewToken("secret", user, time.Hour*24)

	if err != nil {
		return nil, status.Error(codes.Internal, "jwt issue")

	}

	return &authv1.RegisterResponse{Token: token}, nil
}

func (g *GRPCserver) Login(ctx context.Context, in *authv1.LoginRequest) (*authv1.LoginResponse, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	user, err := g.store.User(ctx, in.Email)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "email or password not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(in.Password))
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "email or password not found")
	}
	token, err := jwt.NewToken("secret", user, time.Hour*24)
	if err != nil {
		return nil, status.Error(codes.Internal, "jwt issue")
	}

	return &authv1.LoginResponse{Token: token}, nil
}
