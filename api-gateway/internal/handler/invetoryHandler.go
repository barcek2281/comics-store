package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/barcek2281/comics-store/api-gateway/internal/cache"
	"github.com/barcek2281/comics-store/api-gateway/internal/utils"
	inventoryv1 "github.com/barcek2281/proto/gen/go/inventory"
	"github.com/go-redis/redis/v8"
	"google.golang.org/grpc"
)

const InventoryCachedKey = "inventoryList"

type InventoryHandler struct {
	log             *slog.Logger
	InventoryClient inventoryv1.InventoryClient
	redisClient     *redis.Client
}

func NewInventoryHandler(log *slog.Logger, portInventory int) *InventoryHandler {
	conn, err := grpc.NewClient(fmt.Sprintf("inventory:%d", portInventory), grpc.WithInsecure())
	if err != nil {
		log.Error("failed to connect to inventory service", slog.String("error", err.Error()))
		return nil
	}
	client := inventoryv1.NewInventoryClient(conn)
	return &InventoryHandler{
		log:             log,
		InventoryClient: client,
		redisClient:     cache.NewRedisClient(),
	}
}

func (h *InventoryHandler) Create() http.HandlerFunc {
	type Req struct {
		Title       string `json:"title"`
		Author      string `json:"author"`
		Description string `json:"description"`
		ReleaseDate string `json:"release_date"`
		Price       int64  `json:"price"`
		Quantity    int32  `json:"quantity"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var req Req
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		res, err := h.InventoryClient.Create(ctx, &inventoryv1.CreateRequest{
			Title:       req.Title,
			Author:      req.Author,
			Description: req.Description,
			ReleaseDate: req.ReleaseDate,
			Price:       int64(req.Price),
			Quantity:    int64(req.Quantity),
		})
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to create comic: %v", err), http.StatusInternalServerError)
			return
		}
		ctx, cancel = context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()
		h.redisClient.Del(ctx, InventoryCachedKey)
		slog.Info("updated redis cached data")
		utils.Response(w, r, http.StatusOK, map[string]int64{"id": res.Id})
	}
}

func (h *InventoryHandler) Get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		if id == "" {
			http.Error(w, "missing id parameter", http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		numId, _ := strconv.Atoi(id)

		res, err := h.InventoryClient.Get(ctx, &inventoryv1.GetRequest{Id: int64(numId)})
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to get comic: %v", err), http.StatusInternalServerError)
			return
		}

		utils.Response(w, r, http.StatusOK, res)
	}
}

func (h *InventoryHandler) List() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		page := r.URL.Query().Get("page")
		price := r.URL.Query().Get("price-up")
		cachedKey := InventoryCachedKey

		if price == "" && page == "" {
			ctx, cancel = context.WithTimeout(context.Background(), 60*time.Second)
			defer cancel()
			cached, err := h.redisClient.Get(ctx, cachedKey).Result()
			if err == nil {
				var prod []*inventoryv1.Comics
				_ = json.Unmarshal([]byte(cached), &prod)
				utils.Response(w, r, http.StatusOK, prod)
				slog.Info("get cached data from redis")
				return
			}
		}

		res, err := h.InventoryClient.List(ctx, &inventoryv1.ListRequest{})
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to list comics: %v", err), http.StatusInternalServerError)
			return
		}

		if page != "" {
			pageInt, err := strconv.Atoi(page)
			need_res := []*inventoryv1.Comics{}
			if err == nil {
				for i := range res.Comics {
					if i >= (pageInt-1)*5 && i < (pageInt)*5 {
						need_res = append(need_res, res.Comics[i])
					}
				}
			}
			res.Comics = need_res
		}

		if price != "" {
			priceInt, err := strconv.Atoi(price)
			need_res := []*inventoryv1.Comics{}
			if err == nil {
				for i := range res.Comics {
					if res.Comics[i].Price > float32(priceInt) {
						need_res = append(need_res, res.Comics[i])
					}
				}
			}
			res.Comics = need_res
		}

		if price == "" && page == "" {
			ctx, cancel = context.WithTimeout(context.Background(), 60*time.Second)
			defer cancel()
			b, _ := json.Marshal(res.Comics)
			h.redisClient.Set(ctx, cachedKey, b, time.Minute*5)
			slog.Info("set data for redis")

		}
		utils.Response(w, r, http.StatusOK, res.Comics)
	}
}

func (h *InventoryHandler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		if id == "" {
			http.Error(w, "missing id parameter", http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		numId, _ := strconv.Atoi(id)

		res, err := h.InventoryClient.Delete(ctx, &inventoryv1.DeleteRequest{Id: int64(numId)})
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to delete comic: %v", err), http.StatusInternalServerError)
			return
		}
		ctx, cancel = context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()
		h.redisClient.Del(ctx, InventoryCachedKey)
		slog.Info("updated redis cached data")

		utils.Response(w, r, http.StatusOK, res)
	}
}

func (h *InventoryHandler) Update() http.HandlerFunc {
	type Req struct {
		Id          int64  `json:"id"`
		Title       string `json:"title"`
		Author      string `json:"author"`
		Description string `json:"description"`
		ReleaseDate string `json:"release_date"`
		Price       int64  `json:"price"`
		Quantity    int32  `json:"quantity"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var req Req
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		res, err := h.InventoryClient.Update(ctx, &inventoryv1.UpdateRequest{
			Id:          req.Id,
			Title:       req.Title,
			Author:      req.Author,
			Description: req.Description,
			ReleaseDate: req.ReleaseDate,
			Price:       int64(req.Price),
			Quantity:    int64(req.Quantity),
		})
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to update comic: %v", err), http.StatusInternalServerError)
			return
		}
		ctx, cancel = context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()
		h.redisClient.Del(ctx, InventoryCachedKey)
		slog.Info("updated redis cached data")
		utils.Response(w, r, http.StatusOK, res)
	}
}
