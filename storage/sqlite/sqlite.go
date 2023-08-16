package sqlite

import (
	"context"
	"database/sql"

	"github.com/hahaclassic/golang-telegram-bot.git/lib/errhandling"
	"github.com/hahaclassic/golang-telegram-bot.git/storage"
	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func New(path string) (*Storage, error) {

	db, err := sql.Open("sqlite3", path) // Open database

	if err != nil {
		return nil, errhandling.Wrap("can't open database: %w", err)
	}

	if err := db.Ping(); err != nil { // Check connection to database
		return nil, errhandling.Wrap("can't connect to database", err)
	}

	return &Storage{db: db}, nil
}

// Save() adds page in the storage
func (s *Storage) Save(ctx context.Context, p *storage.Page) error {
	q := `INSERT INTO pages (url, user_name) VALUES (?, ?)`

	if _, err := s.db.ExecContext(ctx, q, p.URL, p.UserName); err != nil {
		return errhandling.Wrap("can't save page", err)
	}

	return nil
}

// PickRandom() picks random page in the storage
func (s *Storage) PickRandom(ctx context.Context, userName string) (*storage.Page, error) {
	q := `SELECT url FROM pages WHERE user_name = ? ORDER BY RANDOM() LIMIT 1`

	var url string

	err := s.db.QueryRowContext(ctx, q, userName).Scan(&url)

	if err == sql.ErrNoRows {
		return nil, storage.ErrNoSavedPages
	}
	if err != nil {
		return nil, errhandling.Wrap("can't pick random page:", err)
	}

	return &storage.Page{
		URL:      url,
		UserName: userName,
	}, nil
}

// Remove() deletes the required page
func (s *Storage) Remove(ctx context.Context, page *storage.Page) error {
	q := `DELETE FROM pages WHERE url = ? AND user_name = ?`

	if _, err := s.db.ExecContext(ctx, q, page.URL, page.UserName); err != nil {
		return errhandling.Wrap("can't remove page", err)
	}

	return nil
}

// IsExists() checks if pages exists in storage
func (s *Storage) IsExist(ctx context.Context, page *storage.Page) (bool, error) {
	q := `SELECT COUNT(*) FROM pages WHERE url = ? AND user_name = ?`

	var count int

	if err := s.db.QueryRowContext(ctx, q, page.URL, page.UserName).Scan(&count); err != nil {
		return false, errhandling.Wrap("can't check if page exists", err)
	}

	return count > 0, nil
}

func (s *Storage) Init(ctx context.Context) error {
	q := `CREATE TABLE IF NOT EXISTS pages (url TEXT, user_name TEXT)`

	_, err := s.db.ExecContext(ctx, q)
	if err != nil {
		return errhandling.Wrap("can't create table", err)
	}

	return nil
}
