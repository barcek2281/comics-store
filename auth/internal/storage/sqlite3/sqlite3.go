package sqlite1488

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/barcek2281/comics-store/auth/internal/model"
	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func NewStorage(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.New"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) Save(user model.User) (int64, error) {

	stmt, err := s.db.Prepare("INSERT INTO users(email, password) VALUES(?, ?)")
	if err != nil {
		return 0, err
	}

	res, err := stmt.Exec(user.Email, user.Password)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *Storage) User(ctx context.Context, email string) (model.User, error) {
	stmt, err := s.db.Prepare("SELECT id, email, password FROM users WHERE email = ?")
	if err != nil {
		return model.User{}, err
	}

	row := stmt.QueryRowContext(ctx, email)

	var user model.User
	err = row.Scan(&user.ID, &user.Email, &user.Password)
	if err != nil {
		return model.User{}, err
	}

	return user, nil
}
