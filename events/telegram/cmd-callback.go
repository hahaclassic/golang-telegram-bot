package telegram

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/hahaclassic/golang-telegram-bot.git/events"
	conc "github.com/hahaclassic/golang-telegram-bot.git/lib/concatenation"
	"github.com/hahaclassic/golang-telegram-bot.git/lib/errhandling"
	"github.com/hahaclassic/golang-telegram-bot.git/storage"
)

func (p *Processor) doCallbackCmd(event *events.Event, meta *CallbackMeta) (err error) {

	defer func() {
		_ = p.tg.AnswerCallbackQuery(meta.QueryID)
		if err == ErrEmptyFolder {
			err = nil
		}
		err = errhandling.WrapIfErr("can't do callback cmd", err)
	}()

	if p.sessions[event.UserID].currentOperation == DeleteLinkCmd {
		p.sessions[event.UserID].tag = strings.TrimSpace(meta.Data)
	} else {
		p.sessions[event.UserID].folderName = strings.TrimSpace(meta.Data)
	}

	switch p.sessions[event.UserID].currentOperation {

	case ChooseFolderForRenamingCmd:
		p.sessions[event.UserID].currentOperation = RenameFolderCmd
		return p.chooseFolderForRenaming(event.ChatID)

	case ChooseLinkForDeletionCmd:
		p.sessions[event.UserID].currentOperation = DeleteLinkCmd
		return p.chooseLinkForDeletion(context.Background(), event)

	case GetNameCmd:
		p.sessions[event.UserID].currentOperation = SaveLinkCmd
		p.sessions[event.UserID].tag = p.sessions[event.UserID].url
		err = p.chooseFolder(context.Background(), event.ChatID, event.UserID)

	case KeyCmd:
		return p.ShowKeys(context.Background(), event.ChatID, event.UserID)

	case SaveLinkCmd:
		p.sessions[event.UserID].status = statusOK
		return p.savePage(context.Background(), event)

	case ShowFolderCmd:
		p.sessions[event.UserID].status = statusOK
		return p.showFolder(context.Background(), event)

	case DeleteFolderCmd:
		p.sessions[event.UserID].status = statusOK
		return p.deleteFolder(context.Background(), event)

	case DeleteLinkCmd:
		p.sessions[event.UserID].status = statusOK
		return p.deleteLink(context.Background(), event)
	}

	return nil
}

func (p *Processor) savePage(ctx context.Context, event *events.Event) (err error) {
	defer func() { err = errhandling.WrapIfErr("can't save page", err) }()

	session := p.sessions[event.UserID]
	folderID, err := p.storage.FolderID(ctx, event.UserID, session.folderName)
	if err != nil {
		return err
	}

	access, err := p.storage.GetAccessLvl(ctx, event.UserID, folderID)
	if err != nil {
		return err
	}
	if access != storage.Owner && access != storage.Editor {
		p.sessions[event.UserID].status = statusOK
		return p.tg.SendMessage(event.ChatID, msgIncorrectAccessLvl)
	}

	page := p.storage.NewPage(session.url, session.tag, folderID)

	isExists, err := p.storage.IsPageExist(ctx, page)
	if err != nil {
		return err
	}
	if isExists {
		return p.tg.SendMessage(event.ChatID, msgAlreadyExists)
	}

	if err := p.storage.SavePage(ctx, page); err != nil {
		return err
	}

	return p.tg.SendMessage(event.ChatID, msgSaved)
}

func (p *Processor) showFolder(ctx context.Context, event *events.Event) (err error) {

	defer func() { err = errhandling.WrapIfErr("can't show folder", err) }()

	folderName := p.sessions[event.UserID].folderName
	folderID, err := p.storage.FolderID(ctx, event.UserID, folderName)
	if err != nil {
		return err
	}
	urls, err := p.storage.GetLinks(ctx, folderID)
	if err != nil {
		return err
	}

	tags, err := p.storage.GetTags(ctx, folderID)
	if err != nil {
		return errhandling.Wrap("can't show folder", err)
	}
	if len(urls) == 0 {
		return p.tg.SendMessage(event.ChatID, msgEmptyFolder)
	}

	result := folderName + ":\n\n" + conc.EnumeratedJoinWithTags(urls, tags)

	return p.tg.SendMessage(event.ChatID, result)
}

func (p *Processor) deleteFolder(ctx context.Context, event *events.Event) error {
	folderName := p.sessions[event.UserID].folderName
	folderID, err := p.storage.FolderID(ctx, event.UserID, folderName)
	if err != nil {
		return err
	}

	access, err := p.storage.GetAccessLvl(ctx, event.UserID, folderID)
	if err != nil {
		return err
	}
	if access != storage.Owner {
		p.sessions[event.UserID].status = statusOK
		return p.tg.SendMessage(event.ChatID, msgIncorrectAccessLvl)
	}

	err = p.storage.RemoveFolder(ctx, folderID)
	if err != nil {
		return errhandling.Wrap("can't delete folder", err)
	}

	return p.tg.SendMessage(event.ChatID, msgFolderDeleted)
}

func (p *Processor) chooseFolderForRenaming(chatID int) error {
	return p.tg.SendMessage(chatID, msgEnterNewFolderName)
}

func (p *Processor) chooseLinkForDeletion(ctx context.Context, event *events.Event) error {

	folderName := p.sessions[event.UserID].folderName
	folderID, err := p.storage.FolderID(ctx, event.UserID, folderName)
	if err != nil {
		return err
	}

	access, err := p.storage.GetAccessLvl(ctx, event.UserID, folderID)
	if err != nil {
		return err
	}
	if access != storage.Owner && access != storage.Editor {
		p.sessions[event.UserID].status = statusOK
		return p.tg.SendMessage(event.ChatID, msgIncorrectAccessLvl)
	}

	urls, err := p.storage.GetTags(ctx, folderID)
	if err != nil {
		return errhandling.Wrap("can't show folder", err)
	}

	if len(urls) == 0 {
		p.tg.SendMessage(event.ChatID, msgEmptyFolder)
		return ErrEmptyFolder
	}

	return p.tg.SendCallbackMessage(event.ChatID, msgChooseLink, urls, urls)
}

func (p *Processor) deleteLink(ctx context.Context, event *events.Event) error {

	folderID, err := p.storage.FolderID(ctx, event.UserID, p.sessions[event.UserID].folderName)
	if err != nil {
		return err
	}

	session := p.sessions[event.UserID]
	// Т.к. поле name является уникальным в отдельной папке, то удаление происходит по нему
	// и URL в следующей строке не имеет значения.
	page := p.storage.NewPage("", session.tag, folderID)
	if page == nil {
		return errors.New("can't delete link: can't create folder")
	}

	err = p.storage.RemovePage(ctx, page)
	if err != nil {
		return err
	}

	return p.tg.SendMessage(event.ChatID, msgPageDeleted)
}

func (p *Processor) ShowKeys(ctx context.Context, ChatID int, UserID int) error {

	folderName := p.sessions[UserID].folderName
	folderID, err := p.storage.FolderID(ctx, UserID, folderName)
	if err != nil {
		return err
	}

	access, err := p.storage.GetAccessLvl(ctx, UserID, folderID)
	if err != nil {
		return err
	}
	if access != storage.Owner {
		p.sessions[UserID].status = statusOK
		return p.tg.SendMessage(ChatID, msgIncorrectAccessLvl)
	}

	keys := []string{}
	names := []string{}
	for lvl := storage.Editor; lvl < storage.Reader; lvl++ {
		key, err := p.storage.GetPassword(ctx, folderID, lvl)
		if err == storage.ErrNoPasswords {
			break
		} else if err != nil {
			return err
		}
		keys = append(keys, key)
		names = append(names, fmt.Sprintf("%s", lvl))
	}

	var message string
	if len(keys) == 0 {
		message = "No passwords"
	} else {
		message = conc.EnumeratedJoinWithTags(keys, names)
	}
	buttons := []string{"Create key", "Delete key"}
	operations := []string{CreateKeyCmd, DeleteKeyCmd}
	return p.tg.SendCallbackMessage(ChatID, message, buttons, operations)
}
