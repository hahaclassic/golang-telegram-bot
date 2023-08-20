package telegram

import (
	"context"
	"errors"
	"log"
	"net/url"
	"strings"

	"github.com/hahaclassic/golang-telegram-bot.git/lib/errhandling"
	"github.com/hahaclassic/golang-telegram-bot.git/storage"
)

func (p *Processor) doCmd(text string, chatID int, userID int) (err error) {

	defer func() {
		if err != nil {
			p.status = statusOK
			p.currentOperation = ""
		}
		if err == ErrNoFolders {
			err = nil
		}
	}()

	text = strings.TrimSpace(text)

	log.Printf("got new command '%s' from '%d'", text, userID)

	if p.status {

		if isAddCmd(text) {
			p.status = statusProcessing
			p.currentOperation = SaveLinkCmd
			p.lastMessage = text
			return p.chooseFolder(context.Background(), chatID, userID)
		}

		switch text {
		case StartCmd:
			return p.sendHello(chatID)
		case RusHelpCmd:
			return p.sendRusHelp(chatID)
		case HelpCmd:
			return p.sendHelp(chatID)
		case RndCmd:
			return p.sendRandom(context.Background(), chatID, userID)

		case ShowFolderCmd:
			p.status = statusProcessing
			p.currentOperation = ShowFolderCmd
			return p.chooseFolder(context.Background(), chatID, userID)

		case CreateFolderCmd:
			p.status = statusProcessing
			p.currentOperation = CreateFolderCmd
			return p.tg.SendMessage(chatID, msgEnterFolderName)

		case ChooseFolderForRenaming:
			p.status = statusProcessing
			p.currentOperation = ChooseFolderForRenaming
			return p.chooseFolder(context.Background(), chatID, userID)

		case DeleteFolderCmd:
			p.status = statusProcessing
			p.currentOperation = DeleteFolderCmd
			return p.chooseFolder(context.Background(), chatID, userID)

		case ChooseLinkForDeletionCmd:
			p.status = statusProcessing
			p.currentOperation = ChooseLinkForDeletionCmd
			return p.chooseFolder(context.Background(), chatID, userID)

		default:
			return p.tg.SendMessage(chatID, msgUnknownCommand)
		}

	} else {
		p.status = statusOK
		switch p.currentOperation {

		case CreateFolderCmd:
			return p.createFolder(context.Background(), chatID, userID, text) // text == folderName

		case RenameFolderCmd:
			return p.renameFolder(context.Background(), chatID, userID, text)

		default:
			return p.tg.SendMessage(chatID, msgUnknownCommand)
		}
	}
}

func (p *Processor) createFolder(ctx context.Context, chatID int, userID int, folder string) (err error) {
	defer func() { err = errhandling.WrapIfErr("can't create folder", err) }()

	ok, err := p.storage.IsFolderExist(ctx, userID, folder)
	if err != nil {
		return err
	}

	if ok {
		p.tg.SendMessage(chatID, msgFolderAlreadyExists)
	} else {
		p.storage.NewFolder(ctx, userID, folder)
		p.tg.SendMessage(chatID, msgNewFolderCreated)
	}

	return nil
}

func (p *Processor) chooseFolder(ctx context.Context, chatID int, userID int) (err error) {
	defer func() {
		if err != ErrNoFolders {
			err = errhandling.WrapIfErr("can't do command: choose folder", err)
		}
	}()

	folders, err := p.storage.GetListOfFolders(ctx, userID)
	if err != nil {
		return err
	}
	if len(folders) == 0 {
		_ = p.tg.SendMessage(chatID, msgNoFolders)
		return ErrNoFolders
	}

	return p.tg.SendCallbackMessage(chatID, msgChooseFolder, folders)
}

func (p *Processor) renameFolder(ctx context.Context, chatID int, userID int, folder string) error {

	ok, err := p.storage.IsFolderExist(ctx, userID, folder)
	if err != nil {
		return errhandling.Wrap("can't rename folder", err)
	}
	if ok {
		return p.tg.SendMessage(chatID, msgCantRename)
	}

	err = p.storage.RenameFolder(ctx, userID, folder, p.lastMessage)
	if err != nil {
		return errhandling.Wrap("can't rename folder", err)
	}

	return p.tg.SendMessage(chatID, msgFolderRenamed)
}

func (p *Processor) sendRandom(ctx context.Context, chatID int, userID int) (err error) {
	defer func() { err = errhandling.WrapIfErr("can't do command: can't send random", err) }()

	page, err := p.storage.PickRandom(ctx, userID)
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

func (p *Processor) sendRusHelp(chatID int) error {
	return p.tg.SendMessage(chatID, msgRusHelp)
}

func (p *Processor) sendHello(chatID int) error {
	return p.tg.SendMessage(chatID, msgHello)
}

func isAddCmd(text string) bool {
	return isURL(text)
}

func isURL(text string) bool {
	// Необходим протокол в ссылке (https://)
	u, err := url.Parse(text)

	return err == nil && u.Host != ""
}
