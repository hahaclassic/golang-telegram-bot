package storage

import (
	"context"
	"errors"
	"strconv"
)

type Storage interface {
	// Pages
	NewPage(url string, tag string, folderID string) *Page
	SavePage(ctx context.Context, p *Page) error
	PickRandom(ctx context.Context, userID int) (string, error)
	RemovePage(ctx context.Context, p *Page) error
	IsPageExist(ctx context.Context, p *Page) (bool, error)

	// Folders
	NewFolder(folderName string, lvl AccessLevel, userID int, username string) *Folder
	IsFolderExist(ctx context.Context, folderID string) (bool, error)
	FolderID(ctx context.Context, userID int, folderName string) (string, error)
	FolderName(ctx context.Context, folderID string) (folderName string, err error)
	AccessLevelByUserID(ctx context.Context, folderID string, userID int) (AccessLevel, error)
	Owner(ctx context.Context, folderID string) (userID int, err error)

	AddFolder(ctx context.Context, folder *Folder) error
	RemoveFolder(ctx context.Context, folderID string) error
	RenameFolder(ctx context.Context, folderID string, folderName string) error
	DeleteAccess(ctx context.Context, userID int, folderID string) error

	Folders(ctx context.Context, userID int) ([][]string, error)
	GetLinks(ctx context.Context, folderID string) ([]string, error)
	GetTags(ctx context.Context, folderID string) ([]string, error)

	// Passwords
	CreatePassword(ctx context.Context, folderID string, accessLvl AccessLevel) error
	GetPassword(ctx context.Context, folderID string, accessLvl AccessLevel) (string, error)
	DeletePassword(ctx context.Context, folderID string, accessLvl AccessLevel) error
	AccessLevelByPassword(ctx context.Context, folderID string, password string) (AccessLevel, error)

	// Crypto (information security)
	// SaveKeys(ctx context.Context, keys [][]string)

	// User Settings
	// ChangeLanguage()
	// OutputFormat()
	// Premium()
}

var (
	ErrNoFolders       = errors.New("No folders")
	ErrNoSavedPages    = errors.New("No saved pages")
	ErrIvalidAccessLvl = errors.New("Invalid access level")
	ErrNoPasswords     = errors.New("No passwords")
	ErrNoRows          = errors.New("Err No Rows")
)

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
	Suspected                          // Last chance to gain access. In case of another refusal, the status will change to banned
	// Those users who have been denied access are marked as blocked. It is necessary in order to protect the user from spam
	Banned
	Undefined
)

func (lvl AccessLevel) String() string {
	return []string{"Owner", "Editor", "Confirmed reader", "Reader", "Banned", "Undefined"}[lvl]
}

func ToAccessLvl(s string) AccessLevel {
	lvl, _ := strconv.Atoi(s)
	for i := Owner; i < Undefined; i++ {
		if int(i) == lvl {
			return i
		}
	}
	return Undefined
}
