package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/barcek2281/comics-store/api-gateway/internal/utils"
	orderv1 "github.com/barcek2281/proto/gen/go/order"
	"google.golang.org/grpc"
)

type OrderHandler struct {
	log         *slog.Logger
	OrderClient orderv1.OrderClient
}

func NewOrderHandler(log *slog.Logger, portOrder int) *OrderHandler {
	conn, err := grpc.Dial(fmt.Sprintf("localhost:%d", portOrder), grpc.WithInsecure())
	if err != nil {
		log.Error("failed to connect to order service", slog.String("error", err.Error()))
		return nil
	}

	client := orderv1.NewOrderClient(conn)
	return &OrderHandler{
		log:         log,
		OrderClient: client,
	}
}

func (h *OrderHandler) CreateOrder() http.HandlerFunc {
	type Item struct {
		ProductId string `json:"product_id"`
		Quantity  int32  `json:"quantity"`
	}
	type Req struct {
		UserId string `json:"user_id"`
		Items  []Item `json:"items"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var req Req
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		var items []*orderv1.OrderItem
		for _, i := range req.Items {
			items = append(items, &orderv1.OrderItem{
				ProductId: i.ProductId,
				Quantity:  i.Quantity,
			})
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		res, err := h.OrderClient.CreateOrder(ctx, &orderv1.CreateOrderRequest{
			UserId: req.UserId,
			Items:  items,
		})
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to create order: %v", err), http.StatusInternalServerError)
			return
		}

		utils.Response(w, r, http.StatusOK, res)

	}
}
func (h *OrderHandler) GetOrder() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orderID := r.URL.Query().Get("id")
		if orderID == "" {
			http.Error(w, "missing order id", http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		res, err := h.OrderClient.GetOrder(ctx, &orderv1.GetOrderRequest{OrderId: orderID})
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to get order: %v", err), http.StatusInternalServerError)
			return
		}

		utils.Response(w, r, http.StatusOK, res)

	}
}
func (h *OrderHandler) UpdateOrder() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orderID := r.URL.Query().Get("id")
		if orderID == "" {
			http.Error(w, "missing order id", http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		res, err := h.OrderClient.UpdateOrder(ctx, &orderv1.GetOrderRequest{OrderId: orderID})
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to update order: %v", err), http.StatusInternalServerError)
			return
		}

		utils.Response(w, r, http.StatusOK, res)

	}
}
func (h *OrderHandler) CloseOrder() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.URL.Query().Get("user_id")
		if userID == "" {
			http.Error(w, "missing user id", http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		res, err := h.OrderClient.CloseOrder(ctx, &orderv1.CloseOrderRequest{UserId: userID})
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to close order: %v", err), http.StatusInternalServerError)
			return
		}
		utils.Response(w, r, http.StatusOK, res)
	}
}
func (h *OrderHandler) DeleteOrder() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.URL.Query().Get("user_id")
		if userID == "" {
			http.Error(w, "missing user id", http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		res, err := h.OrderClient.DeleteOrder(ctx, &orderv1.DeleteOrderRequest{UserId: userID})
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to delete order: %v", err), http.StatusInternalServerError)
			return
		}

		utils.Response(w, r, http.StatusOK, res)
	}
}
func (h *OrderHandler) ListOrders() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.URL.Query().Get("user_id")
		if userID == "" {
			http.Error(w, "missing user id", http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		res, err := h.OrderClient.ListOrders(ctx, &orderv1.OrderListRequest{UserId: userID})
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to list orders: %v", err), http.StatusInternalServerError)
			return
		}

		utils.Response(w, r, http.StatusOK, res)
	}
}
