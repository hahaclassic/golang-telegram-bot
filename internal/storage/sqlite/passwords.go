package sqlite

import (
	"context"
	"database/sql"
	"strings"

	"github.com/google/uuid"
	"github.com/hahaclassic/golang-telegram-bot.git/internal/storage"
	"github.com/hahaclassic/golang-telegram-bot.git/lib/errhandling"
)

func (s *Storage) CreatePassword(ctx context.Context, folderID string, accessLvl storage.AccessLevel) error {

	var q string

	unic, err := uuid.NewRandom()
	if err != nil {
		return errhandling.WrapIfErr("error while generating password", err)
	}
	password := strings.ToUpper(unic.String()[:8])

	_, err = s.GetPassword(ctx, folderID, accessLvl)

	if err == storage.ErrNoPasswords {
		q = `INSERT INTO passwords (folder_id, access_level, password) VALUES (?, ?, ?)`
		_, err = s.db.ExecContext(ctx, q, folderID, accessLvl, password)
	} else if err == nil {
		q = `UPDATE passwords SET password = ? WHERE folder_id = ? AND access_level = ?`
		_, err = s.db.ExecContext(ctx, q, password, folderID, accessLvl)
	}

	return errhandling.WrapIfErr("can't save password", err)
}

func (s *Storage) GetPassword(ctx context.Context, folderID string, accessLvl storage.AccessLevel) (string, error) {
	var password string

	q := `SELECT password FROM passwords WHERE folder_id = ? AND access_level = ?`

	err := s.db.QueryRowContext(ctx, q, folderID, accessLvl).Scan(&password)
	if err == sql.ErrNoRows {
		return "", storage.ErrNoPasswords
	}
	if err != nil {
		return "", errhandling.Wrap("cant get password", err)
	}

	return "KEY" + folderID + password, nil
}

func (s *Storage) DeletePassword(ctx context.Context, folderID string, accessLvl storage.AccessLevel) error {

	q := `DELETE FROM passwords WHERE folder_id = ? AND access_level = ?`

	_, err := s.db.ExecContext(ctx, q, folderID, accessLvl)
	if err == sql.ErrNoRows {
		return storage.ErrNoPasswords
	}

	return errhandling.WrapIfErr("can' remove page", err)
}

// GetAccessLvl returns the user's access level to the specified folder
func (s *Storage) AccessLevelByPassword(ctx context.Context, folderID string, password string) (storage.AccessLevel, error) {

	var accessLvl storage.AccessLevel

	q := `SELECT access_level FROM passwords WHERE folder_id = ? AND password = ?`

	err := s.db.QueryRowContext(ctx, q, folderID, password).Scan(&accessLvl)

	if err == sql.ErrNoRows {
		return storage.Undefined, nil
	}
	if err != nil {
		return storage.Undefined, errhandling.Wrap("cant get access_level", err)
	}

	return accessLvl, nil
}
