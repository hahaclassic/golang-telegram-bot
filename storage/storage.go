package storage

import (
	"context"
	"errors"
)

type Storage interface {
	Save(ctx context.Context, p *Page) error
	PickRandom(ctx context.Context, userName string) (*Page, error)
	Remove(ctx context.Context, p *Page) error

	GetFolder(ctx context.Context, page *Page) ([]string, error)
	RemoveFolder(ctx context.Context, page *Page) error

	IsExist(ctx context.Context, p *Page) (bool, error)
	IsFolderExist(ctx context.Context, page *Page) (bool, error)
}

var ErrNoSavedPages = errors.New("no saved pages")

type Page struct {
	URL      string
	UserName string
	Folder   string
}
