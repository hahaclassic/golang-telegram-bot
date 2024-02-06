package storage

import (
	"context"
	"errors"
)

type Storage interface {
	// Pages
	NewPage(url string, name string, folderID string) *Page
	SavePage(ctx context.Context, p *Page) error
	PickRandom(ctx context.Context, userID int) (string, error)
	RemovePage(ctx context.Context, p *Page) error
	IsPageExist(ctx context.Context, p *Page) (bool, error)

	// Folders
	NewFolder(folderName string, lvl AccessLevel, userID int, username string) *Folder
	IsFolderExist(ctx context.Context, folderID string) (bool, error)
	// FolderID(ctx context.Context, userID int, folderName string) (string, error)
	GetAccessLvl(ctx context.Context, userID int, folderID string) (AccessLevel, error)
	AddFolder(ctx context.Context, folder *Folder) error
	RemoveFolder(ctx context.Context, folderID string) error
	RenameFolder(ctx context.Context, folderID string, folderName string) error
	DeleteAccess(ctx context.Context, userID int, folderID string) error

	GetFolders(ctx context.Context, userID int) ([][]string, error)
	GetLinks(ctx context.Context, folderID string) ([]string, error)
	GetTags(ctx context.Context, folderID string) ([]string, error)

	// Passwords
	CreatePassword(ctx context.Context, folderID string, accessLvl AccessLevel) error
	GetPassword(ctx context.Context, folderID string, accessLvl AccessLevel) (string, error)
	DeletePassword(ctx context.Context, folderID string, accessLvl AccessLevel) error

	// Crypto (information security)
	// SaveKeys(ctx context.Context, keys [][]string)

	// User Settings
	// ChangeLanguage()
	// OutputFormat()
	// Premium()
}

var ErrNoSavedPages = errors.New("No saved pages")
var ErrIvalidAccessLvl = errors.New("Invalid access level")

type Page struct {
	URL      string
	Tag      string
	FolderID string
}

type Folder struct {
	ID        string
	Name      string
	AccessLvl AccessLevel
	UserID    int
	Username  string
}

type AccessLevel int

const (
	Owner           AccessLevel = iota // Owner: All possible actions with the folder and its contents are available
	Editor                             // Editor: add/delete links (always after confirmation)
	ConfirmedReader                    // reading only after confirmation
	Reader                             // reading only
	// Those users who have been denied access are marked as blocked. It is necessary in order to protect the user from spam
	Banned
	Undefined
)
