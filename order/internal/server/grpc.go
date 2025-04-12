package server

import (
	"context"
	"time"

	"github.com/barcek2281/comics-store/order/internal/storage"
	orderv1 "github.com/barcek2281/proto/gen/go/order"
	"github.com/google/uuid"
)

type GRPCserver struct {
	orderv1.UnimplementedOrderServer
	store *storage.Storage
}

func NewGRPCserver(store *storage.Storage) *GRPCserver {
	return &GRPCserver{
		store: store,
	}
}

func (g *GRPCserver) CreateOrder(ctx context.Context, in *orderv1.CreateOrderRequest) (*orderv1.CreateOrderResponse, error) {
	orderID := uuid.New().String()
	totalPrice := float32(0)

	for _, item := range in.Items {
		totalPrice += float32(item.Quantity) * 10.0
	}

	order := &orderv1.Order{
		Id:         orderID,
		UserId:     in.UserId,
		Items:      in.Items,
		TotalPrice: totalPrice,
		Status:     "created",
		CreatedAt:  time.Now().Format(time.RFC3339),
	}

	err := g.store.CreateOrder(ctx, order)
	if err != nil {
		return nil, err
	}

	return &orderv1.CreateOrderResponse{
		OrderId: orderID,
		Status:  "created",
	}, nil
}

func (g *GRPCserver) GetOrder(ctx context.Context, in *orderv1.GetOrderRequest) (*orderv1.Order, error) {
	order, err := g.store.GetOrderByID(ctx, in.OrderId)
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (g *GRPCserver) UpdateOrder(ctx context.Context, in *orderv1.GetOrderRequest) (*orderv1.UpdateOrderResponce, error) {
	// For simplicity, let’s just mark the order as “updated”
	err := g.store.UpdateOrderStatus(ctx, in.OrderId, "updated")
	if err != nil {
		return &orderv1.UpdateOrderResponce{Status: "failed"}, err
	}
	return &orderv1.UpdateOrderResponce{Status: "updated"}, nil
}

func (g *GRPCserver) CloseOrder(ctx context.Context, in *orderv1.CloseOrderRequest) (*orderv1.CloseOrderResponce, error) {
	err := g.store.CloseOrderByUserID(ctx, in.UserId)
	if err != nil {
		return &orderv1.CloseOrderResponce{
			IsChanged: false,
			Status:    "failed",
		}, err
	}
	return &orderv1.CloseOrderResponce{
		IsChanged: true,
		Status:    "closed",
	}, nil
}

func (g *GRPCserver) DeleteOrder(ctx context.Context, in *orderv1.DeleteOrderRequest) (*orderv1.CloseOrderResponce, error) {
	err := g.store.DeleteOrderByUserID(ctx, in.UserId)
	if err != nil {
		return &orderv1.CloseOrderResponce{
			IsChanged: false,
			Status:    "failed",
		}, err
	}
	return &orderv1.CloseOrderResponce{
		IsChanged: true,
		Status:    "deleted",
	}, nil
}

func (g *GRPCserver) ListOrders(ctx context.Context, in *orderv1.OrderListRequest) (*orderv1.OrderListResponse, error) {
	orders, err := g.store.ListOrdersByUserID(ctx, in.UserId)
	if err != nil {
		return nil, err
	}
	return &orderv1.OrderListResponse{Orders: orders}, nil
}
