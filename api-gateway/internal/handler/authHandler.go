package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	authv1 "github.com/barcek2281/proto-comics/gen/go/auth"
	"google.golang.org/grpc"
)

type AuthHandler struct {
	log        *slog.Logger
	AuthClient authv1.AuthClient
}

func NewAuthHandler(log *slog.Logger, portAuth int) *AuthHandler {
	conn, err := grpc.Dial(fmt.Sprintf("localhost:%d", portAuth), grpc.WithInsecure())
	if err != nil {
		return nil
	}
	AuthClient := authv1.NewAuthClient(conn)
	return &AuthHandler{
		log:        log,
		AuthClient: AuthClient,
	}
}

func (h *AuthHandler) Register() http.HandlerFunc {
	type Req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		var req Req
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {

			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
		defer cancel()

		res, err := h.AuthClient.Register(ctx, &authv1.RegisterRequest{
			Email: req.Email,
			Password: req.Password,
		})
		if err != nil {

			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"token": res.Token})
	}
}
