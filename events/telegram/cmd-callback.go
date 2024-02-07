package telegram

import (
	"context"
	"strings"

	"github.com/hahaclassic/golang-telegram-bot.git/events"
	conc "github.com/hahaclassic/golang-telegram-bot.git/lib/concatenation"
	"github.com/hahaclassic/golang-telegram-bot.git/lib/errhandling"
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

	// case ChooseFolderForRenamingCmd:
	// 	p.sessions[meta.UserID].currentOperation = RenameFolderCmd
	// 	return p.chooseFolderForRenaming(meta.ChatID)

	// case ChooseLinkForDeletionCmd:
	// 	p.sessions[meta.UserID].currentOperation = DeleteLinkCmd
	// 	return p.chooseLinkForDeletion(context.Background(), meta)

	case GetNameCmd:
		p.sessions[event.UserID].currentOperation = SaveLinkCmd
		p.sessions[event.UserID].tag = p.sessions[event.UserID].url
		err = p.chooseFolder(context.Background(), event.ChatID, event.UserID)

		// case SaveLinkCmd:
		// 	p.sessions[meta.UserID].status = statusOK
		// 	return p.savePage(context.Background(), meta)

		// case ShowFolderCmd:
		// 	p.sessions[meta.UserID].status = statusOK
		// 	return p.showFolder(context.Background(), meta)

		// case DeleteFolderCmd:
		// 	p.sessions[meta.UserID].status = statusOK
		// 	return p.deleteFolder(context.Background(), meta)

		// case DeleteLinkCmd:
		// 	p.sessions[meta.UserID].status = statusOK
		// 	return p.deleteLink(context.Background(), meta)
	}

	return nil
}

// func (p *Processor) savePage(ctx context.Context, meta *CallbackMeta) (err error) {
// 	defer func() { err = errhandling.WrapIfErr("can't save page", err) }()

// 	session := p.sessions[meta.UserID]
// 	page := p.storage.NewPage(session.url, session.name, meta.UserID, session.folder)

// 	isExists, err := p.storage.IsExist(ctx, page)
// 	if err != nil {
// 		return err
// 	}
// 	if isExists {
// 		return p.tg.SendMessage(meta.ChatID, msgAlreadyExists)
// 	}

// 	if err := p.storage.Save(ctx, page); err != nil {
// 		return err
// 	}

// 	if err := p.tg.SendMessage(meta.ChatID, msgSaved); err != nil {
// 		return err
// 	}

// 	return nil
// }

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

// func (p *Processor) deleteFolder(ctx context.Context, meta *CallbackMeta) error {
// 	folder := p.sessions[meta.UserID].folder
// 	err := p.storage.RemoveFolder(ctx, meta.UserID, folder)
// 	if err != nil {
// 		return errhandling.Wrap("can't delete folder", err)
// 	}

// 	return p.tg.SendMessage(meta.ChatID, msgFolderDeleted)
// }

// func (p *Processor) chooseFolderForRenaming(chatID int) error {
// 	return p.tg.SendMessage(chatID, msgEnterNewFolderName)
// }

// func (p *Processor) chooseLinkForDeletion(ctx context.Context, meta *CallbackMeta) error {

// 	folder := p.sessions[meta.UserID].folder
// 	urls, err := p.storage.GetNames(ctx, meta.UserID, folder)
// 	if err != nil {
// 		return errhandling.Wrap("can't show folder", err)
// 	}

// 	if len(urls) == 0 {
// 		p.tg.SendMessage(meta.ChatID, msgEmptyFolder)
// 		return ErrEmptyFolder
// 	}

// 	return p.tg.SendCallbackMessage(meta.ChatID, msgChooseLink, urls)
// }

// func (p *Processor) deleteLink(ctx context.Context, meta *CallbackMeta) error {

// 	session := p.sessions[meta.UserID]
// 	// Т.к. поле name является уникальным в отдельной папке, то удаление происходит по нему
// 	// и URL в следующей строке не имеет значения.
// 	page := p.storage.NewPage("", session.name, meta.UserID, session.folder)

// 	err := p.storage.Remove(ctx, page)
// 	if err != nil {
// 		return err
// 	}

// 	return p.tg.SendMessage(meta.ChatID, msgPageDeleted)
// }
