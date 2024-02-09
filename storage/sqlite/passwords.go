package sqlite

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/hahaclassic/golang-telegram-bot.git/lib/errhandling"
	"github.com/hahaclassic/golang-telegram-bot.git/storage"
)

func (s *Storage) CreatePassword(ctx context.Context, folderID string, accessLvl storage.AccessLevel) error {

	q := `INSERT INTO passwords (folder_id, access_level, password) VALUES (?, ?, ?)`

	unic, err := uuid.NewRandom()
	if err != nil {
		return errhandling.WrapIfErr("error while generating password", err)
	}
	password := unic.String()[:8]

	_, err = s.db.ExecContext(ctx, q, folderID, accessLvl, password)

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

	return password, nil
}

func (s *Storage) DeletePassword(ctx context.Context, folderID string, accessLvl storage.AccessLevel) error {

	q := `DELETE FROM passwords WHERE folder_id = ? AND access_level = ?`

	_, err := s.db.ExecContext(ctx, q, folderID, accessLvl)

	return errhandling.WrapIfErr("can' remove page", err)
}
