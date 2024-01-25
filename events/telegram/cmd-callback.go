package telegram

import (
	"context"
	"strings"

	conc "github.com/hahaclassic/golang-telegram-bot.git/lib/concatenation"
	"github.com/hahaclassic/golang-telegram-bot.git/lib/errhandling"
)

func (p *Processor) doCallbackCmd(text string, meta *CallbackMeta) (err error) {
	defer func() {
		switch p.sessions[meta.UserID].currentOperation {
		case ChooseFolderForRenaming:
			p.changeSessionData(meta.UserID, &Session{text, RenameFolderCmd, statusProcessing})
		case ChooseLinkForDeletionCmd:
			p.changeSessionData(meta.UserID, &Session{text, DeleteLinkCmd, statusProcessing})
		default:
			p.changeSessionData(meta.UserID, &Session{text, "", statusOK})
		}

		_ = p.tg.AnswerCallbackQuery(meta.QueryID)
		if err == ErrEmptyFolder {
			err = nil
		}
		err = errhandling.WrapIfErr("can't do callback cmd", err)
	}()

	text = strings.TrimSpace(text)

	switch p.sessions[meta.UserID].currentOperation {
	case SaveLinkCmd:
		return p.savePage(context.Background(), meta, text)

	case ShowFolderCmd:
		return p.showFolder(context.Background(), meta, text)

	case ChooseFolderForRenaming:
		return p.chooseFolderForRenaming(meta.ChatID)

	case DeleteFolderCmd:
		return p.deleteFolder(context.Background(), meta, text)

	case ChooseLinkForDeletionCmd:
		return p.chooseLinkForDeletion(context.Background(), meta, text)

	case DeleteLinkCmd:
		return p.deleteLink(context.Background(), meta, text)
	}

	return nil
}

func (p *Processor) savePage(ctx context.Context, meta *CallbackMeta, folder string) (err error) {
	defer func() { err = errhandling.WrapIfErr("can't save page", err) }()

	page := p.storage.NewPage(p.sessions[meta.UserID].lastMessage, meta.UserID, folder)

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

func (p *Processor) showFolder(ctx context.Context, meta *CallbackMeta, folder string) error {

	urls, err := p.storage.GetFolder(ctx, meta.UserID, folder)
	if err != nil {
		return errhandling.Wrap("can't show folder", err)
	}

	if len(urls) == 0 {
		return p.tg.SendMessage(meta.ChatID, msgEmptyFolder)
	}

	result := folder + ":\n" + conc.EnumeratedJoin(urls)

	return p.tg.SendMessage(meta.ChatID, result)
}

func (p *Processor) deleteFolder(ctx context.Context, meta *CallbackMeta, folder string) error {

	err := p.storage.RemoveFolder(ctx, meta.UserID, folder)
	if err != nil {
		return errhandling.Wrap("can't delete folder", err)
	}

	return p.tg.SendMessage(meta.ChatID, msgFolderDeleted)
}

func (p *Processor) chooseFolderForRenaming(chatID int) error {
	return p.tg.SendMessage(chatID, msgEnterNewFolderName)
}

func (p *Processor) chooseLinkForDeletion(ctx context.Context, meta *CallbackMeta, folder string) error {

	urls, err := p.storage.GetFolder(ctx, meta.UserID, folder)
	if err != nil {
		return errhandling.Wrap("can't show folder", err)
	}

	if len(urls) == 0 {
		p.tg.SendMessage(meta.ChatID, msgEmptyFolder)
		return ErrEmptyFolder
	}

	return p.tg.SendCallbackMessage(meta.ChatID, msgChooseLink, urls)
}

func (p *Processor) deleteLink(ctx context.Context, meta *CallbackMeta, link string) error {

	page := p.storage.NewPage(link, meta.UserID, p.sessions[meta.UserID].lastMessage)

	err := p.storage.Remove(ctx, page)
	if err != nil {
		return err
	}

	return p.tg.SendMessage(meta.ChatID, msgPageDeleted)
}
