package telegram

import (
	"context"
	"strings"

	"github.com/hahaclassic/golang-telegram-bot.git/lib/errhandling"
)

func (p *Processor) doCallbackCmd(text string, meta *CallbackMeta) (err error) {
	defer func() {
		_ = p.tg.AnswerCallbackQuery(meta.QueryID)
		if err == ErrEmptyFolder {
			err = nil
		}
		err = errhandling.WrapIfErr("can't do callback cmd", err)
	}()

	text = strings.TrimSpace(text)
	if p.sessions[meta.UserID].currentOperation == DeleteLinkCmd {
		p.sessions[meta.UserID].name = text
	} else {
		p.sessions[meta.UserID].folder = text
	}

	switch p.sessions[meta.UserID].currentOperation {

	case ChooseFolderForRenamingCmd:
		p.sessions[meta.UserID].currentOperation = RenameFolderCmd
		return p.chooseFolderForRenaming(meta.ChatID)

	case ChooseLinkForDeletionCmd:
		p.sessions[meta.UserID].currentOperation = DeleteLinkCmd
		return p.chooseLinkForDeletion(context.Background(), meta)

	case GetNameCmd:
		p.sessions[meta.UserID].currentOperation = RenameFolderCmd
		p.sessions[meta.UserID].name = trimLink(p.sessions[meta.UserID].url)
		err = p.chooseFolder(context.Background(), meta.ChatID, meta.UserID)

	case SaveLinkCmd:
		p.sessions[meta.UserID].status = statusOK
		return p.savePage(context.Background(), meta)

	case ShowFolderCmd:
		p.sessions[meta.UserID].status = statusOK
		return p.showFolder(context.Background(), meta)

	case DeleteFolderCmd:
		p.sessions[meta.UserID].status = statusOK
		return p.deleteFolder(context.Background(), meta)

	case DeleteLinkCmd:
		p.sessions[meta.UserID].status = statusOK
		return p.deleteLink(context.Background(), meta)
	}

	return nil
}

func (p *Processor) savePage(ctx context.Context, meta *CallbackMeta) (err error) {
	defer func() { err = errhandling.WrapIfErr("can't save page", err) }()

	session := p.sessions[meta.UserID]
	page := p.storage.NewPage(session.url, session.name, meta.UserID, session.folder)

	isExists, err := p.storage.IsExist(ctx, page)
	if err != nil {
		return err
	}
	if isExists {
		return p.tg.SendMessage(meta.ChatID, msgAlreadyExists)
	}

	if err := p.storage.Save(ctx, page); err != nil {
		return err
	}

	if err := p.tg.SendMessage(meta.ChatID, msgSaved); err != nil {
		return err
	}

	return nil
}

func (p *Processor) showFolder(ctx context.Context, meta *CallbackMeta) error {

	folder := p.sessions[meta.UserID].folder
	urls, err := p.storage.GetLinks(ctx, meta.UserID, folder)
	if err != nil {
		return errhandling.Wrap("can't show folder", err)
	}

	names, err := p.storage.GetNames(ctx, meta.UserID, folder)
	if err != nil {
		return errhandling.Wrap("can't show folder", err)
	}

	if len(urls) == 0 {
		return p.tg.SendMessage(meta.ChatID, msgEmptyFolder)
	}

	result := folder + ":\n\n" + linkList(urls, names)

	return p.tg.SendMessage(meta.ChatID, result)
}

func (p *Processor) deleteFolder(ctx context.Context, meta *CallbackMeta) error {
	folder := p.sessions[meta.UserID].folder
	err := p.storage.RemoveFolder(ctx, meta.UserID, folder)
	if err != nil {
		return errhandling.Wrap("can't delete folder", err)
	}

	return p.tg.SendMessage(meta.ChatID, msgFolderDeleted)
}

func (p *Processor) chooseFolderForRenaming(chatID int) error {
	return p.tg.SendMessage(chatID, msgEnterNewFolderName)
}

func (p *Processor) chooseLinkForDeletion(ctx context.Context, meta *CallbackMeta) error {

	folder := p.sessions[meta.UserID].folder
	urls, err := p.storage.GetNames(ctx, meta.UserID, folder)
	if err != nil {
		return errhandling.Wrap("can't show folder", err)
	}

	if len(urls) == 0 {
		p.tg.SendMessage(meta.ChatID, msgEmptyFolder)
		return ErrEmptyFolder
	}

	return p.tg.SendCallbackMessage(meta.ChatID, msgChooseLink, urls)
}

func (p *Processor) deleteLink(ctx context.Context, meta *CallbackMeta) error {

	session := p.sessions[meta.UserID]
	// Т.к. поле name является уникальным в отдельной папке, то удаление происходит по нему
	// и URL в следующей строке не имеет значения.
	page := p.storage.NewPage("", session.name, meta.UserID, session.folder)

	err := p.storage.Remove(ctx, page)
	if err != nil {
		return err
	}

	return p.tg.SendMessage(meta.ChatID, msgPageDeleted)
}
