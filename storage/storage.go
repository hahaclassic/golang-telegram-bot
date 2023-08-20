package storage

import (
	"context"
	"errors"
)

type Storage interface {
	Save(ctx context.Context, p *Page) error
	PickRandom(ctx context.Context, userID int) (*Page, error)
	Remove(ctx context.Context, p *Page) error
	IsExist(ctx context.Context, p *Page) (bool, error)

	NewFolder(ctx context.Context, userID int, folder string) error
	RemoveFolder(ctx context.Context, userID int, folder string) error
	GetFolder(ctx context.Context, userID int, folder string) ([]string, error)
	GetListOfFolders(ctx context.Context, userID int) (names []string, err error)
	IsFolderExist(ctx context.Context, userID int, folder string) (bool, error)
	RenameFolder(ctx context.Context, userID int, newFolder, oldFolder string) error
}

var ErrNoSavedPages = errors.New("no saved pages")

type Page struct {
	URL    string
	UserID int
	Folder string
}
