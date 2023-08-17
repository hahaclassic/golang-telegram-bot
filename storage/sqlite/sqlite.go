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
	q := `INSERT INTO pages (url, user_name, folder) VALUES (?, ?, ?)`

	if _, err := s.db.ExecContext(ctx, q, p.URL, p.UserName, p.Folder); err != nil {
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
	q := `DELETE FROM pages WHERE url = ? AND user_name = ? AND folder = ?`

	if _, err := s.db.ExecContext(ctx, q, page.URL, page.UserName, page.Folder); err != nil {
		return errhandling.Wrap("can't remove page", err)
	}

	return nil
}

func (s *Storage) RemoveFolder(ctx context.Context, page *storage.Page) error {
	q := `DELETE FROM pages WHERE user_name = ? AND folder = ?`

	if _, err := s.db.ExecContext(ctx, q, page.UserName, page.Folder); err != nil {
		return errhandling.Wrap("can't remove folder", err)
	}

	return nil
}

func (s *Storage) GetFolder(ctx context.Context, page *storage.Page) (urls []string, err error) {
	defer func() { err = errhandling.WrapIfErr("can't get folder", err) }()

	q := `SELECT url FROM pages WHERE user_name = ? AND folder = ?`

	rows, err := s.db.QueryContext(ctx, q, page.UserName, page.Folder)
	if err != nil {
		return nil, err
	}

	var temp string

	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&temp); err != nil {
			return nil, err
		}
		urls = append(urls, temp)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return urls, nil
}

// func (s *Storage) AllFolders(ctx context.Context, page *storage.Page) (err error) {
// 	defer func() { err = errhandling.WrapIfErr("can't check folders", err) }()

// 	// q := `SELECT COUNT(*) FROM pages WHERE user_name = ?` // Check count of all pages

// 	// var foldersCount int

// 	// if err := s.db.QueryRowContext(ctx, q, page.UserName).Scan(&foldersCount); err != nil {
// 	// 	return err
// 	// }

// 	folders = make([]string)
// 	q = `SELECT DISTINCT folder FROM pages WHERE user_name = ?` // Get all pages

// 	if err := s.db.QueryContext(ctx, q, page.UserName).Scan(folders); err != nil {
// 		return nil, err
// 	}

// 	return folders, nil
// }

func (s *Storage) IsFolderExist(ctx context.Context, page *storage.Page) (bool, error) {
	q := `SELECT COUNT(*) FROM pages WHERE user_name = ? AND folder = ?`

	var count int

	if err := s.db.QueryRowContext(ctx, q, page.UserName, page.Folder).Scan(&count); err != nil {
		return false, errhandling.Wrap("can't check if page exists", err)
	}

	return count > 0, nil
}

// IsExists() checks if pages exists in storage
func (s *Storage) IsExist(ctx context.Context, page *storage.Page) (bool, error) {
	q := `SELECT COUNT(*) FROM pages WHERE url = ? AND user_name = ? AND folder = ?`

	var count int

	if err := s.db.QueryRowContext(ctx, q, page.URL, page.UserName, page.Folder).Scan(&count); err != nil {
		return false, errhandling.Wrap("can't check if page exists", err)
	}

	return count > 0, nil
}

func (s *Storage) Init(ctx context.Context) error {
	q := `CREATE TABLE IF NOT EXISTS pages (url TEXT, user_name TEXT, folder TEXT)`

	_, err := s.db.ExecContext(ctx, q)
	if err != nil {
		return errhandling.Wrap("can't create table", err)
	}

	return nil
}
