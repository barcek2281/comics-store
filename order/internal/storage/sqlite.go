package storage

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"strconv"

	inventoryv1 "github.com/barcek2281/proto/gen/go/inventory"
	orderv1 "github.com/barcek2281/proto/gen/go/order"

	_ "github.com/mattn/go-sqlite3"
	"google.golang.org/grpc"
)

type Storage struct {
	db              *sql.DB
	InventoryCLient inventoryv1.InventoryClient
}

func NewStorage(storagePath string) (*Storage, error) {
	conn, err := grpc.NewClient(fmt.Sprintf("inventory:%s", "50052"), grpc.WithInsecure())
	if err != nil {
		return nil, nil
	}
	client := inventoryv1.NewInventoryClient(conn)

	const op = "storage.sqlite.New"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{
		db: db,
		InventoryCLient: client,
		}, nil
}

func (s *Storage) CreateOrder(ctx context.Context, order *orderv1.Order) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	for _, item := range order.Items {
		n, _ := strconv.Atoi(item.ProductId)
		comics, err := s.InventoryCLient.Get(ctx, &inventoryv1.GetRequest{Id: int64(n)})
		if err != nil {
			tx.Rollback()
			return err
		}
		slog.Info("info", "item quantity", item.Quantity, "avaible", comics.Quantity)
		if item.Quantity > comics.Quantity {
			tx.Rollback()
			slog.Error("cannot fit with items", "item quantity", item.Quantity, "avaible", comics.Quantity)

			return fmt.Errorf("not enough stock for comic id %s", item.ProductId)
		}
	}

	_, err = tx.ExecContext(ctx,
		`INSERT INTO orders (id, user_id, total_price, status, created_at)
		 VALUES (?, ?, ?, ?, ?)`,
		order.Id, order.UserId, order.TotalPrice, order.Status, order.CreatedAt,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	for _, item := range order.Items {
		_, err := tx.ExecContext(ctx,
			`INSERT INTO order_items (order_id, product_id, quantity)
			 VALUES (?, ?, ?)`,
			order.Id, item.ProductId, item.Quantity,
		)
		if err != nil {
			tx.Rollback()
			slog.Error("error write to db", "error", err)

			return err
		}
	}

	return tx.Commit()
}

func (s *Storage) GetOrderByID(ctx context.Context, orderID string) (*orderv1.Order, error) {
	order := &orderv1.Order{}
	err := s.db.QueryRowContext(ctx,
		`SELECT id, user_id, total_price, status, created_at FROM orders WHERE id = ?`,
		orderID,
	).Scan(&order.Id, &order.UserId, &order.TotalPrice, &order.Status, &order.CreatedAt)
	if err != nil {
		return nil, err
	}

	rows, err := s.db.QueryContext(ctx,
		`SELECT product_id, quantity FROM order_items WHERE order_id = ?`,
		orderID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		item := &orderv1.OrderItem{}
		if err := rows.Scan(&item.ProductId, &item.Quantity); err != nil {
			return nil, err
		}
		order.Items = append(order.Items, item)
	}

	return order, nil
}

func (s *Storage) UpdateOrderStatus(ctx context.Context, orderID string, status string) error {
	_, err := s.db.ExecContext(ctx,
		`UPDATE orders SET status = ? WHERE id = ?`,
		status, orderID,
	)
	return err
}

func (s *Storage) CloseOrderByUserID(ctx context.Context, userID string) error {
	_, err := s.db.ExecContext(ctx,
		`UPDATE orders SET status = 'closed' WHERE user_id = ?`,
		userID,
	)
	return err
}

func (s *Storage) DeleteOrderByUserID(ctx context.Context, userID string) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	rows, err := tx.QueryContext(ctx, `SELECT id FROM orders WHERE user_id = ?`, userID)
	if err != nil {
		tx.Rollback()
		slog.Error("error to find", "error", err)

		return err
	}
	defer rows.Close()

	var orderIDs []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			tx.Rollback()
			return err
		}
		orderIDs = append(orderIDs, id)
	}

	for _, id := range orderIDs {
		_, err = tx.ExecContext(ctx, `DELETE FROM order_items WHERE order_id = ?`, id)
		if err != nil {
			tx.Rollback()
			slog.Error("error to delete", "error", err)

			return err
		}
	}

	_, err = tx.ExecContext(ctx, `DELETE FROM orders WHERE user_id = ?`, userID)
	if err != nil {
		tx.Rollback()
		slog.Error("error to delete", "error", err)

		return err
	}

	return tx.Commit()
}
func (s *Storage) ListOrdersByUserID(ctx context.Context, userID string) ([]*orderv1.Order, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, total_price, status, created_at FROM orders WHERE user_id = ?`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*orderv1.Order
	for rows.Next() {
		order := &orderv1.Order{UserId: userID}
		err := rows.Scan(&order.Id, &order.TotalPrice, &order.Status, &order.CreatedAt)
		if err != nil {
			return nil, err
		}

		itemRows, err := s.db.QueryContext(ctx,
			`SELECT product_id, quantity FROM order_items WHERE order_id = ?`,
			order.Id,
		)
		if err != nil {
			return nil, err
		}

		for itemRows.Next() {
			item := &orderv1.OrderItem{}
			if err := itemRows.Scan(&item.ProductId, &item.Quantity); err != nil {
				itemRows.Close()
				return nil, err
			}
			order.Items = append(order.Items, item)
		}
		itemRows.Close()

		orders = append(orders, order)
	}

	return orders, nil
}
