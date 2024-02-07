package sqlite

import (
	"context"
	"database/sql"

	"github.com/hahaclassic/golang-telegram-bot.git/lib/errhandling"
	"github.com/hahaclassic/golang-telegram-bot.git/storage"
)

func (s *Storage) NewPage(url string, tag string, folderID string) *storage.Page {
	return &storage.Page{
		URL:      url,
		Tag:      tag,
		FolderID: folderID,
	}
}

// Save() adds page in the storage
// UNSAFE. The access level is not checked.
func (s *Storage) SavePage(ctx context.Context, p *storage.Page) error {
	q := `INSERT INTO pages (url, tag, folder_id) VALUES (?, ?, ?)`

	_, err := s.db.ExecContext(ctx, q, p.URL, p.Tag, p.FolderID)

	return errhandling.WrapIfErr("can't save page", err)
}

// PickRandom() picks random page in the storage
func (s *Storage) PickRandom(ctx context.Context, userID int) (string, error) {

	var url, folderID string

	q := `SELECT folder_id FROM folders WHERE user_id = ? AND access_level != ? ORDER BY RANDOM() LIMIT 1`

	err := s.db.QueryRowContext(ctx, q, userID, storage.Banned).Scan(&folderID)
	if err == sql.ErrNoRows {
		return "", storage.ErrNoSavedPages
	}
	if err != nil {
		return "", errhandling.Wrap("can't pick random page:", err)
	}

	q = `SELECT url FROM pages WHERE folder_id = ? ORDER BY RANDOM() LIMIT 1`

	err = s.db.QueryRowContext(ctx, q, folderID).Scan(&url)
	if err == sql.ErrNoRows {
		return "", storage.ErrNoSavedPages
	}
	if err != nil {
		return "", errhandling.Wrap("can't pick random page:", err)
	}

	return url, nil
}

// Remove() deletes the required page
// UNSAFE. The access level is not checked.
func (s *Storage) RemovePage(ctx context.Context, page *storage.Page) error {
	q := `DELETE FROM pages WHERE (url = ? OR tag = ?) AND folder_id = ?`

	_, err := s.db.ExecContext(ctx, q, page.URL, page.Tag, page.FolderID)

	return errhandling.WrapIfErr("can' remove page", err)
}

// IsExists() checks if pages exists in storage
func (s *Storage) IsPageExist(ctx context.Context, page *storage.Page) (bool, error) {
	q := `SELECT COUNT(*) FROM pages WHERE (url = ? OR tag = ?) AND folder_id = ?`

	var count int

	if err := s.db.QueryRowContext(ctx, q, page.URL, page.Tag, page.FolderID).Scan(&count); err != nil {
		return false, errhandling.Wrap("can't check if page exists", err)
	}

	return count > 0, nil
}
