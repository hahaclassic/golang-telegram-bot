package storage

import (
	"context"
	"errors"
)

type Storage interface {
	Save(ctx context.Context, p *Page) error
	PickRandom(ctx context.Context, userName string) (*Page, error)
	Remove(ctx context.Context, p *Page) error
	IsExist(ctx context.Context, p *Page) (bool, error)

	NewFolder(ctx context.Context, username string, folder string) error
	RemoveFolder(ctx context.Context, username string, folder string) error
	GetFolder(ctx context.Context, username string, folder string) ([]string, error)
	GetListOfFolders(ctx context.Context, username string) (names []string, err error)
	IsFolderExist(ctx context.Context, username string, folder string) (bool, error)
	RenameFolder(ctx context.Context, username, newFolder, oldFolder string) error
}

var ErrNoSavedPages = errors.New("no saved pages")

type Page struct {
	URL      string
	UserName string
	Folder   string
}
