package telegram

import (
	"context"
	"errors"
	"log"
	"net/url"
	"strings"

	conc "github.com/hahaclassic/golang-telegram-bot.git/lib/concatenation"
	"github.com/hahaclassic/golang-telegram-bot.git/lib/errhandling"
	"github.com/hahaclassic/golang-telegram-bot.git/storage"
)

const (
	RndCmd          = "/rnd"
	HelpCmd         = "/help"
	StartCmd        = "/start"
	DeleteLinkCmd   = "/delete_link"   // Удаляет ссылку из нужной папки
	ShowFolderCmd   = "/folder"        // Показывает содержимое папки
	CreateFolderCmd = "/create"        // Создает новую папку
	DeleteFolderCmd = "/delete_folder" // Удаляет папку
	ChangeFolderCmd = "/change"        // Меняет местонахождение папки
	RenameFolderCmd = "/rename"        // Изменяет название папки
	SaveLink        = "/save"
)

func (p *Processor) doCmd(text string, chatID int, username string) error {
	text = strings.TrimSpace(text)

	log.Printf("got new command '%s' from '%s'", text, username)

	if p.status {

		if isAddCmd(text) {
			p.status = statusProcessing
			p.lastMessage = text
			p.currentOperation = SaveLink
			return p.tg.SendMessage(chatID, msgEnterFolderName)
		}

		switch text {
		case RndCmd:
			return p.sendRandom(context.Background(), chatID, username)
		case HelpCmd:
			return p.sendHelp(chatID)
		case StartCmd:
			return p.sendHello(chatID)
		case DeleteFolderCmd:
			p.status = statusProcessing
			p.currentOperation = DeleteFolderCmd
			return p.tg.SendMessage(chatID, msgEnterFolderName)
		case CreateFolderCmd:
			p.status = statusProcessing
			p.currentOperation = CreateFolderCmd
			return p.tg.SendMessage(chatID, msgEnterFolderName)
		case ShowFolderCmd:
			p.status = statusProcessing
			p.currentOperation = ShowFolderCmd
			return p.tg.SendMessage(chatID, msgEnterFolderName)
		default:
			return p.tg.SendMessage(chatID, msgUnknownCommand)
		}

	} else {

		p.status = statusOK
		switch p.currentOperation {
		case SaveLink:
			return p.savePage(context.Background(), chatID, p.lastMessage, username, text) // text == folderName

		// case CreateFolderCmd:
		// 	return p.createFolder(context.Background(), chatID, username, text) // text == folderName
		// }

		case ShowFolderCmd:
			return p.showFolder(context.Background(), chatID, username, text)
		default:
			return p.tg.SendMessage(chatID, msgUnknownCommand)
		}
	}
}

func (p *Processor) showFolder(ctx context.Context, chatID int, username string, folder string) (err error) {
	defer func() { err = errhandling.WrapIfErr("can't do command: show folder", err) }()

	page := &storage.Page{
		UserName: username,
		Folder:   folder,
	}

	isExists, err := p.storage.IsFolderExist(ctx, page)
	if err != nil {
		return err
	}
	if !isExists {
		return p.tg.SendMessage(chatID, msgFolderNotExists)
	}

	urls, err := p.storage.GetFolder(ctx, page)
	if err != nil {
		return err
	}

	resultMessage := "folder " + folder + ":\n" + conc.EnumeratedJoin(urls)

	return p.tg.SendMessage(chatID, resultMessage)
}

func (p *Processor) savePage(ctx context.Context, chatID int, pageURL string, username string, folder string) (err error) {
	defer func() { err = errhandling.WrapIfErr("can't do command: save page", err) }()

	page := &storage.Page{
		URL:      pageURL,
		UserName: username,
		Folder:   folder,
	}

	isExists, err := p.storage.IsExist(ctx, page)
	if err != nil {
		return err
	}
	if isExists {
		return p.tg.SendMessage(chatID, msgAlreadyExists)
	}

	if err := p.storage.Save(ctx, page); err != nil {
		return err
	}

	if err := p.tg.SendMessage(chatID, msgSaved); err != nil {
		return err
	}

	return nil
}

func (p *Processor) sendRandom(ctx context.Context, chatID int, username string) (err error) {
	defer func() { err = errhandling.WrapIfErr("can't do command: can't send random", err) }()

	page, err := p.storage.PickRandom(ctx, username)
	if err != nil && !errors.Is(err, storage.ErrNoSavedPages) {
		return err
	}
	if errors.Is(err, storage.ErrNoSavedPages) {
		return p.tg.SendMessage(chatID, msgNoSavedPages)
	}

	if err := p.tg.SendMessage(chatID, page.URL); err != nil {
		return err
	}

	return p.storage.Remove(ctx, page)
}

func (p *Processor) sendHelp(chatID int) error {
	return p.tg.SendMessage(chatID, msgHelp)
}

func (p *Processor) sendHello(chatID int) error {
	return p.tg.SendMessage(chatID, msgHello)
}

func isAddCmd(text string) bool {
	return isURL(text)
}

func isURL(text string) bool {
	// Необходим протокол в ссылке (https:/)
	u, err := url.Parse(text)

	return err == nil && u.Host != ""
}
