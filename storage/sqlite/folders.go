package sqlite

import (
	"context"

	"github.com/hahaclassic/golang-telegram-bot.git/lib/errhandling"
)

// NewFolder creates a new folder for user in the storage
func (s *Storage) NewFolder(ctx context.Context, userID int, folder string) error {

	q := `INSERT INTO folders (userID, folder) VALUES (?, ?)`

	if _, err := s.db.ExecContext(ctx, q, userID, folder); err != nil {
		return errhandling.Wrap("can't save page", err)
	}

	return nil
}

// RemoveFolder() deletes the required folder
func (s *Storage) RemoveFolder(ctx context.Context, userID int, folder string) error {
	q := `DELETE FROM pages WHERE userID = ? AND folder = ?`

	if _, err := s.db.ExecContext(ctx, q, userID, folder); err != nil {
		return errhandling.Wrap("can't remove folder from table 'pages'", err)
	}

	q = `DELETE FROM folders WHERE userID = ? AND folder = ?`

	if _, err := s.db.ExecContext(ctx, q, userID, folder); err != nil {
		return errhandling.Wrap("can't remove folder from table 'folders'", err)
	}

	return nil
}

// GetFolder() returns list of URL links in folder
func (s *Storage) GetFolder(ctx context.Context, userID int, folder string) (urls []string, err error) {
	defer func() { err = errhandling.WrapIfErr("can't get folder", err) }()

	q := `SELECT url FROM pages WHERE userID = ? AND folder = ?`

	rows, err := s.db.QueryContext(ctx, q, userID, folder)
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

// GetListOfFolders() get list of folders in the storage
func (s *Storage) GetListOfFolders(ctx context.Context, userID int) (names []string, err error) {
	defer func() { err = errhandling.WrapIfErr("can't select all folders", err) }()

	q := `SELECT folder FROM folders WHERE userID = ?` // Get all folders

	rows, err := s.db.QueryContext(ctx, q, userID)
	if err != nil {
		return nil, err
	}

	var temp string

	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&temp); err != nil {
			return nil, err
		}
		names = append(names, temp)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return names, nil
}

// IsFolderExists() checks if folder exists in the storage
func (s *Storage) IsFolderExist(ctx context.Context, userID int, folder string) (bool, error) {
	q := `SELECT COUNT(*) FROM folders WHERE userID = ? AND folder = ?`

	var count int

	if err := s.db.QueryRowContext(ctx, q, userID, folder).Scan(&count); err != nil {
		return false, errhandling.Wrap("can't check if page exists", err)
	}

	return count > 0, nil
}

// RenameFolder() changes the folder name to a new one
func (s *Storage) RenameFolder(ctx context.Context, userID int, newFolder, oldFolder string) error {
	q := `UPDATE folders SET folder = ? WHERE userID = ? AND folder = ?`

	if _, err := s.db.ExecContext(ctx, q, newFolder, userID, oldFolder); err != nil {
		return errhandling.Wrap("can't rename folder", err)
	}

	q = `UPDATE pages SET folder = ? WHERE userID = ? AND folder = ?`

	if _, err := s.db.ExecContext(ctx, q, newFolder, userID, oldFolder); err != nil {
		return errhandling.Wrap("can't rename folder", err)
	}

	return nil
}
