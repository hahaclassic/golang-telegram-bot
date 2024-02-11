package sqlite

import (
	"context"
	"database/sql"
	"strings"

	"github.com/google/uuid"
	"github.com/hahaclassic/golang-telegram-bot.git/lib/errhandling"
	"github.com/hahaclassic/golang-telegram-bot.git/storage"
)

// NewFolder creates a new folder for user in the storage
func (s *Storage) NewFolder(folderName string, lvl storage.AccessLevel, userID int, username string) *storage.Folder {

	id, err := uuid.NewRandom()
	if err != nil {
		return nil
	}
	folderID := strings.ToUpper(id.String())
	folderID = folderID[len(folderID)-12:]

	return &storage.Folder{
		ID:        folderID,
		Name:      folderName,
		AccessLvl: lvl,
		UserID:    userID,
		Username:  username,
	}
}

// IsFolderExists() checks if folder exists in the storage
func (s *Storage) IsFolderExist(ctx context.Context, folderID string) (bool, error) {
	q := `SELECT COUNT(*) FROM folders WHERE folder_id = ?`

	var count int

	if err := s.db.QueryRowContext(ctx, q, folderID).Scan(&count); err != nil {
		return false, errhandling.Wrap("can't check if folder exists", err)
	}

	return count > 0, nil
}

// FolderID()
func (s *Storage) FolderID(ctx context.Context, userID int, folderName string) (folderID string, err error) {

	q := `SELECT folder_id FROM folders WHERE user_id = ? AND folder_name = ?`

	err = s.db.QueryRowContext(ctx, q, userID, folderName).Scan(&folderID)
	if err == sql.ErrNoRows {
		return "", storage.ErrNoFolders
	}
	if err != nil {
		return "", err
	}

	return folderID, nil
}

// FolderID()
func (s *Storage) FolderName(ctx context.Context, folderID string) (folderName string, err error) {

	q := `SELECT folder_name FROM folders WHERE folder_id = ?`

	err = s.db.QueryRowContext(ctx, q, folderID).Scan(&folderName)
	if err == sql.ErrNoRows {
		return "", storage.ErrNoFolders
	}
	if err != nil {
		return "", err
	}

	return folderName, nil
}

// GetAccessLvl returns the user's access level to the specified folder
func (s *Storage) AccessLevelByUserID(ctx context.Context, folderID string, userID int) (storage.AccessLevel, error) {

	var accessLvl storage.AccessLevel

	q := `SELECT access_level FROM folders WHERE folder_id = ? AND user_id = ?`

	if err := s.db.QueryRowContext(ctx, q, folderID, userID).Scan(&accessLvl); err != nil {
		return storage.Undefined, errhandling.Wrap("cant get access_level", err)
	}

	return accessLvl, nil
}

func (s *Storage) Owner(ctx context.Context, folderID string) (userID int, err error) {

	q := `SELECT user_id FROM folders WHERE folder_id = ? AND access_level = ?`

	if err := s.db.QueryRowContext(ctx, q, folderID, storage.Owner).Scan(&userID); err != nil {
		return 0, errhandling.Wrap("cant get folder's owner", err)
	}

	return userID, nil
}

// AddFolder() adds a record that the user has access to the folder
// UNSAFE. The existence of the folder is not checked.
func (s *Storage) AddFolder(ctx context.Context, folder *storage.Folder) error {

	q := `INSERT INTO folders (folder_id, folder_name, access_level, 
		user_id, username) VALUES (?, ?, ?, ?, ?)`

	_, err := s.db.ExecContext(ctx, q, folder.ID, folder.Name, folder.AccessLvl,
		folder.UserID, folder.Username)

	return errhandling.WrapIfErr("can't add folder", err)
}

// RemoveFolder() deletes the required folder
// UNSAFE. The access level is not checked.
func (s *Storage) RemoveFolder(ctx context.Context, folderID string) error {

	q := `DELETE FROM folders WHERE folder_id = ?`

	if _, err := s.db.ExecContext(ctx, q, folderID); err != nil {
		return errhandling.Wrap("can't remove folder from table 'folders'", err)
	}

	q = `DELETE FROM pages WHERE folder_id = ?`

	if _, err := s.db.ExecContext(ctx, q, folderID); err != nil {
		return errhandling.Wrap("can't remove folder's pages from table 'pages'", err)
	}

	q = `DELETE FROM passwords WHERE folder_id = ?`

	if _, err := s.db.ExecContext(ctx, q, folderID); err != nil {
		return errhandling.Wrap("can't remove folder's passwords from table 'passwords'", err)
	}

	return nil
}

// DeleteAccess() deletes the user's access to the folder, but does not delete the folder itself
func (s *Storage) DeleteAccess(ctx context.Context, userID int, folderID string) error {
	q := `DELETE FROM folders WHERE folder_id = ? AND user_id = ?`

	_, err := s.db.ExecContext(ctx, q, folderID, userID)
	if err == sql.ErrNoRows {
		return storage.ErrNoFolders
	}

	return nil
}

// RenameFolder() changes the folder name to a new one
func (s *Storage) RenameFolder(ctx context.Context, folderID string, folderName string) error {
	q := `UPDATE folders SET folder_name = ? WHERE folder_id = ?`

	_, err := s.db.ExecContext(ctx, q, folderName, folderID)

	return errhandling.WrapIfErr("can't rename folder", err)
}

// Folders() returns folders' names and IDs in the [][]string where
// folders[0] - folderID
// folders[1] - folderName
func (s *Storage) Folders(ctx context.Context, userID int) (folders [][]string, err error) {
	defer func() { err = errhandling.WrapIfErr("can't get folders", err) }()

	foldersIDs, err := s.FoldersIDs(ctx, userID)
	if err != nil {
		return nil, err
	}
	foldersNames, err := s.FoldersNames(ctx, userID)
	if err != nil {
		return nil, err
	}

	folders = append(folders, foldersIDs)
	folders = append(folders, foldersNames)
	return folders, nil
}

func (s *Storage) FoldersIDs(ctx context.Context, userID int) (ids []string, err error) {
	q := `SELECT folder_id FROM folders WHERE user_id = ? AND access_level <= ?`

	rows, err := s.db.QueryContext(ctx, q, userID, storage.Reader)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var temp string
	for rows.Next() {
		if err := rows.Scan(&temp); err != nil {
			return nil, err
		}
		ids = append(ids, temp)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return ids, nil
}

func (s *Storage) FoldersNames(ctx context.Context, userID int) (names []string, err error) {
	q := `SELECT folder_name FROM folders WHERE user_id = ? AND access_level <= ?`

	rows, err := s.db.QueryContext(ctx, q, userID, storage.Reader)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var temp string

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

// GetLinks() returns list of URL links in folder
func (s *Storage) GetLinks(ctx context.Context, folderID string) (urls []string, err error) {
	defer func() { err = errhandling.WrapIfErr("can't get folder", err) }()

	q := `SELECT url FROM pages WHERE folder_id = ?`

	rows, err := s.db.QueryContext(ctx, q, folderID)
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

// GetLinks() returns list of URL names (tags) in folder
func (s *Storage) GetTags(ctx context.Context, folderID string) (tags []string, err error) {
	defer func() { err = errhandling.WrapIfErr("can't get tags", err) }()

	q := `SELECT tag FROM pages WHERE folder_id = ?`

	rows, err := s.db.QueryContext(ctx, q, folderID)
	if err != nil {
		return nil, err
	}

	var temp string

	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&temp); err != nil {
			return nil, err
		}
		tags = append(tags, temp)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tags, nil
}
