package sqlite

import (
	"context"
	"database/sql"

	"github.com/hahaclassic/golang-telegram-bot.git/lib/errhandling"
	"github.com/hahaclassic/golang-telegram-bot.git/storage"
)

func (s *Storage) NewPage(url string, name string, userID int, folder string) *storage.Page {
	return &storage.Page{
		URL:    url,
		Name:   name,
		UserID: userID,
		Folder: folder,
	}
}

// Save() adds page in the storage
func (s *Storage) Save(ctx context.Context, p *storage.Page) error {
	q := `INSERT INTO pages (url, name, userID, folder) VALUES (?, ?, ?, ?)`

	if _, err := s.db.ExecContext(ctx, q, p.URL, p.Name, p.UserID, p.Folder); err != nil {
		return errhandling.Wrap("can't save page", err)
	}

	return nil
}

// PickRandom() picks random page in the storage
func (s *Storage) PickRandom(ctx context.Context, userID int) (*storage.Page, error) {
	q := `SELECT url FROM pages WHERE userID = ? ORDER BY RANDOM() LIMIT 1`

	var url string

	err := s.db.QueryRowContext(ctx, q, userID).Scan(&url)

	if err == sql.ErrNoRows {
		return nil, storage.ErrNoSavedPages
	}
	if err != nil {
		return nil, errhandling.Wrap("can't pick random page:", err)
	}

	return &storage.Page{
		URL:    url,
		UserID: userID,
	}, nil
}

// Remove() deletes the required page
func (s *Storage) Remove(ctx context.Context, page *storage.Page) error {
	q := `DELETE FROM pages WHERE name = ? AND userID = ? AND folder = ?`

	if _, err := s.db.ExecContext(ctx, q, page.Name, page.UserID, page.Folder); err != nil {
		return errhandling.Wrap("can't remove page", err)
	}

	return nil
}

// IsExists() checks if pages exists in storage
func (s *Storage) IsExist(ctx context.Context, page *storage.Page) (bool, error) {
	q := `SELECT COUNT(*) FROM pages WHERE (url = ? OR name = ?) AND userID = ? AND folder = ?`

	var count int

	if err := s.db.QueryRowContext(ctx, q, page.URL, page.Name, page.UserID, page.Folder).Scan(&count); err != nil {
		return false, errhandling.Wrap("can't check if page exists", err)
	}

	return count > 0, nil
}
