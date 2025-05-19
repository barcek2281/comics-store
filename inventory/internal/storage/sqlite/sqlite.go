package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/barcek2281/comics-store/inventory/internal/model"
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

func (s *Storage) Create(comics model.Comics) (int64, error) {
	stmt, err := s.db.Prepare(`
		INSERT INTO comics(title, author, description, release_date, price, quantity)
		VALUES (?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return 0, fmt.Errorf("prepare insert: %w", err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(
		comics.Title,
		comics.Author,
		comics.Description,
		comics.ReleaseDate,
		comics.Price,
		comics.Quantity,
	)
	if err != nil {
		return 0, fmt.Errorf("exec insert: %w", err)
	}

	return res.LastInsertId()
}
// Delete deletes a comic by ID
func (s *Storage) Delete(id int64) error {
	_, err := s.db.Exec("DELETE FROM comics WHERE id = ?", id)
	return err
}

// Get fetches a comic by ID
func (s *Storage) Get(id int64) (model.Comics, error) {
	row := s.db.QueryRow("SELECT id, title, author, description, release_date, price, quantity FROM comics WHERE id = ?", id)

	var comic model.Comics
	err := row.Scan(
		&comic.ID,
		&comic.Title,
		&comic.Author,
		&comic.Description,
		&comic.ReleaseDate,
		&comic.Price,
		&comic.Quantity,
	)
	if err != nil {
		return model.Comics{}, err
	}
	return comic, nil
}

// List fetches all comics
func (s *Storage) List() ([]model.Comics, error) {
	rows, err := s.db.Query("SELECT id, title, author, description, release_date, price, quantity FROM comics")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comics []model.Comics
	for rows.Next() {
		var comic model.Comics
		err := rows.Scan(
			&comic.ID,
			&comic.Title,
			&comic.Author,
			&comic.Description,
			&comic.ReleaseDate,
			&comic.Price,
			&comic.Quantity,
		)
		if err != nil {
			return nil, err
		}
		comics = append(comics, comic)
	}
	return comics, nil
}

// Update modifies an existing comic
func (s *Storage) Update(comic model.Comics) error {
	ctx, cn := context.WithTimeout(context.Background(), time.Second*60)
	defer cn()
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	_, err = tx.Exec(`
		UPDATE comics
		SET title = ?, author = ?, description = ?, release_date = ?, price = ?, quantity = ?
		WHERE id = ?
	`,
		comic.Title,
		comic.Author,
		comic.Description,
		comic.ReleaseDate,
		comic.Price,
		comic.Quantity,
		comic.ID,
	)
	if err != nil {
		tx.Rollback()
	}
	return tx.Commit()
}
