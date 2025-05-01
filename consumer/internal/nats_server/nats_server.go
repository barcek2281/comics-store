package nats_server

import (
	"consumer/internal/models"
	"consumer/internal/store/sqlite"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"strconv"

	inventoryv1 "github.com/barcek2281/proto/gen/go/inventory"
	"github.com/nats-io/nats.go"
	"google.golang.org/grpc"
)

type NatsServer struct {
	NC              *nats.Conn
	inventoryClient inventoryv1.InventoryClient
	store           *sqlite.Store
}

func NewNatsServer(store *sqlite.Store, inventoryPort int) *NatsServer {
	nc, err := nats.Connect("nats://nats:4222")
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	conn, err := grpc.NewClient(fmt.Sprintf("inventory:%d", inventoryPort), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("error to connect grpc, error: %v, port: %d", err, inventoryPort)
	}
	client := inventoryv1.NewInventoryClient(conn)
	return &NatsServer{
		NC:              nc,
		store:           store,
		inventoryClient: client,
	}
}

func (n *NatsServer) HandleCreateOrder(m *nats.Msg) {
	var order models.Order

	if err := json.Unmarshal(m.Data, &order); err != nil {
		slog.Error("cannot parse data", "error", err)
		return
	}

	err := n.store.WriteCreatedOrder(order)
	if err != nil {
		slog.Error("error to write db", "error", err)
	}
	ctx := context.Background()
	for _, item := range order.Items {

		newId, _ := strconv.Atoi(item.ProductId)
		comics, err := n.inventoryClient.Get(ctx, &inventoryv1.GetRequest{
			Id: int64(newId),
		})

		if err != nil {
			slog.Warn("error to get value")
			continue
		}

		n.inventoryClient.Update(ctx, &inventoryv1.UpdateRequest{
			Id:          int64(newId),
			Title:       comics.Title,
			Quantity:    int64(comics.Quantity - item.Quantity),
			Author:      comics.Author,
			Description: comics.Description,
			Price:       int64(comics.Price),
		})
	}
	slog.Info("recieve data", "data", order)
}
