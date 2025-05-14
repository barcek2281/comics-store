package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/barcek2281/comics-store/api-gateway/internal/cache"
	"github.com/barcek2281/comics-store/api-gateway/internal/utils"
	orderv1 "github.com/barcek2281/proto/gen/go/order"
	"github.com/go-redis/redis/v8"
	"google.golang.org/grpc"
)

type OrderHandler struct {
	log         *slog.Logger
	OrderClient orderv1.OrderClient
	redisClient *redis.Client
}

func NewOrderHandler(log *slog.Logger, portOrder int) *OrderHandler {
	conn, err := grpc.Dial(fmt.Sprintf("order:%d", portOrder), grpc.WithInsecure())
	if err != nil {
		log.Error("failed to connect to order service", slog.String("error", err.Error()))
		return nil
	}

	client := orderv1.NewOrderClient(conn)
	return &OrderHandler{
		log:         log,
		OrderClient: client,
		redisClient: cache.NewRedisClient(),
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
			utils.Error(w, r, http.StatusBadRequest, fmt.Errorf("invalid request body"))
			h.log.Error("invalid body")
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
			utils.Error(w, r, http.StatusInternalServerError, fmt.Errorf("failed to create order: %v", err))
			return
		}
		cachedKey := fmt.Sprintf("listOrders:%s", req.UserId)
		ctx, cancel = context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()
		h.redisClient.Del(ctx, cachedKey)
		slog.Info("updated redis cached data")

		utils.Response(w, r, http.StatusOK, res)

	}
}
func (h *OrderHandler) GetOrder() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orderID := r.URL.Query().Get("id")
		if orderID == "" {
			utils.Error(w, r, http.StatusBadRequest, fmt.Errorf("missing order id"))
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		res, err := h.OrderClient.GetOrder(ctx, &orderv1.GetOrderRequest{OrderId: orderID})
		if err != nil {
			utils.Error(w, r, http.StatusInternalServerError, fmt.Errorf("failed to get order: %v", err))
			return
		}

		utils.Response(w, r, http.StatusOK, res)

	}
}
func (h *OrderHandler) UpdateOrder() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orderID := r.URL.Query().Get("id")
		if orderID == "" {
			utils.Error(w, r, http.StatusBadRequest, fmt.Errorf("missing order id"))
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		res, err := h.OrderClient.UpdateOrder(ctx, &orderv1.GetOrderRequest{OrderId: orderID})
		if err != nil {
			utils.Error(w, r, http.StatusInternalServerError, fmt.Errorf("failed to update order: %v", err))
			return
		}

		cachedKey := fmt.Sprintf("listOrders:%s", orderID)
		ctx, cancel = context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()
		h.redisClient.Del(ctx, cachedKey)
		slog.Info("updated redis cached data")
		utils.Response(w, r, http.StatusOK, res)

	}
}
func (h *OrderHandler) CloseOrder() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.URL.Query().Get("user_id")
		if userID == "" {
			utils.Error(w, r, http.StatusBadRequest, fmt.Errorf("missing user id"))
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		res, err := h.OrderClient.CloseOrder(ctx, &orderv1.CloseOrderRequest{UserId: userID})
		if err != nil {
			utils.Error(w, r, http.StatusInternalServerError, fmt.Errorf("failed to close order: %v", err))
			return
		}

		cachedKey := fmt.Sprintf("listOrders:%s", userID)
		ctx, cancel = context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()
		h.redisClient.Del(ctx, cachedKey)
		slog.Info("updated redis cached data")

		utils.Response(w, r, http.StatusOK, res)
	}
}
func (h *OrderHandler) DeleteOrder() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.URL.Query().Get("user_id")
		if userID == "" {
			utils.Error(w, r, http.StatusBadRequest, fmt.Errorf("missing user id"))
			return
		}
		cachedKey := fmt.Sprintf("listOrders:%s", userID)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		res, err := h.OrderClient.DeleteOrder(ctx, &orderv1.DeleteOrderRequest{UserId: userID})
		if err != nil {
			utils.Error(w, r, http.StatusInternalServerError, fmt.Errorf("failed to delete order: %v", err))
			return
		}

		ctx, cancel = context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()
		h.redisClient.Del(ctx, cachedKey)
		slog.Info("updated redis cached data")

		utils.Response(w, r, http.StatusOK, res)
	}
}
func (h *OrderHandler) ListOrders() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.URL.Query().Get("user_id")
		if userID == "" {
			utils.Error(w, r, http.StatusBadRequest, fmt.Errorf("missing user id"))
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		cachedKey := fmt.Sprintf("listOrders:%s", userID)

		data, err := h.redisClient.Get(ctx, cachedKey).Result()
		if err == nil {
			var or orderv1.OrderListResponse

			_ = json.Unmarshal([]byte(data), &or)
			utils.Response(w, r, http.StatusOK, or)
			slog.Info("get cached data from redis")
			return
		}

		ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		res, err := h.OrderClient.ListOrders(ctx, &orderv1.OrderListRequest{UserId: userID})
		if err != nil {
			utils.Error(w, r, http.StatusInternalServerError, fmt.Errorf("failed to list orders: %v", err))
			return
		}

		ctx, cancel = context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()
		b, _ := json.Marshal(res)
		h.redisClient.Set(ctx, cachedKey, b, time.Minute*5)
		slog.Info("set data for redis")

		utils.Response(w, r, http.StatusOK, res)
	}
}
