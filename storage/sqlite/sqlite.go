package sqlite

import (
	"context"
	"database/sql"

	"github.com/hahaclassic/golang-telegram-bot.git/lib/errhandling"
	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

// New() create a new database
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

// Init() create tables in the database
func (s *Storage) Init(ctx context.Context) error {
	q := `CREATE TABLE IF NOT EXISTS pages (url TEXT, name TEXT, userID INTEGER, folder TEXT)`

	_, err := s.db.ExecContext(ctx, q)
	if err != nil {
		return errhandling.Wrap("can't create table 'pages'", err)
	}

	q = `CREATE TABLE IF NOT EXISTS folders (userID INTEGER, folder TEXT)`
	_, err = s.db.ExecContext(ctx, q)
	if err != nil {
		return errhandling.Wrap("can't create table 'folders", err)
	}

	return nil
}
