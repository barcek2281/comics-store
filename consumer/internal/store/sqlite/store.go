package sqlite

import (
	"consumer/internal/models"
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type Store struct {
	db *sql.DB
}

func NewStore(storagePath string) *Store {
	const op = "storage.sqlite.New"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		log.Fatalf("cannot connect to db, error: %v", err)
	}

	return &Store{
		db: db,
	}
}

func (s *Store) WriteCreatedOrder(order models.Order) error {

	_, err := s.db.Exec("INSERT INTO order_log(id, price, status, user_id, create_at) VALUES (?, ?, ?, ?, ?)", order.Id, order.TotalPrice, order.Status, order.UserId, order.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}
