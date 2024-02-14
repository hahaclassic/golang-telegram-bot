package telegram

import (
	"context"
	"errors"

	"github.com/hahaclassic/golang-telegram-bot.git/events"
	conc "github.com/hahaclassic/golang-telegram-bot.git/lib/concatenation"
	"github.com/hahaclassic/golang-telegram-bot.git/lib/errhandling"
	"github.com/hahaclassic/golang-telegram-bot.git/storage"
)

// Done
func (p *Processor) chooseFolder(ctx context.Context, chatID int, userID int) (err error) {
	defer func() {
		err = errhandling.WrapIfErr("can't do command: chooseFolder()", err)
	}()

	folders, err := p.storage.Folders(ctx, userID)
	if err != nil {
		return err
	}
	if len(folders[0]) == 0 {
		p.sessions[userID].status = statusOK
		return p.tg.SendMessage(chatID, msgNoFolders)
	}

	messageID, err := p.tg.SendCallbackMessage(chatID, msgChooseFolder, folders[1], folders[0])
	if err == nil {
		p.sessions[userID].lastMessageID = messageID
	}

	return err
}

func (p *Processor) showFolder(ctx context.Context, ChatID int, UserID int) (err error) {

	defer func() { err = errhandling.WrapIfErr("can't show folder", err) }()

	session := p.sessions[UserID]
	urls, err := p.storage.GetLinks(ctx, session.folderID)
	if err != nil {
		return err
	}

	tags, err := p.storage.GetTags(ctx, session.folderID)
	if err != nil {
		return errhandling.Wrap("can't show folder", err)
	}
	if len(urls) == 0 {
		return p.tg.SendMessage(ChatID, msgEmptyFolder)
	}

	folderName, err := p.storage.FolderName(ctx, session.folderID)
	if err != nil {
		return err
	}
	result := folderName + ":\n\n" + conc.EnumeratedJoinWithTags(urls, tags)

	return p.tg.SendMessage(ChatID, result)
}

// event.Text == folderName
func (p *Processor) createFolder(ctx context.Context, event *events.Event) (err error) {
	defer func() { err = errhandling.WrapIfErr("can't create folder", err) }()

	_, err = p.storage.FolderID(ctx, event.UserID, event.Text)
	if err == nil {
		return p.tg.SendMessage(event.ChatID, msgFolderAlreadyExists)
	}

	var folder *storage.Folder

	var i int
	for i = 0; i < maxAttemts; i++ {
		folder = p.storage.NewFolder(event.Text, storage.Owner, event.UserID, event.Username)

		ok, err := p.storage.IsFolderExist(ctx, folder.ID)
		if err == nil && !ok {
			break
		}
	}
	if i == 100 {
		return errors.New("can't create unic folderID")
	}

	err = p.storage.AddFolder(ctx, folder)
	if err != nil {
		return err
	}

	return p.tg.SendMessage(event.ChatID, msgNewFolderCreated)
}

// event.Text == newFolderName
func (p *Processor) renameFolder(ctx context.Context, event *events.Event) (err error) {

	defer func() { err = errhandling.WrapIfErr("can't rename folder", err) }()

	access, err := p.storage.AccessLevelByUserID(ctx, p.sessions[event.UserID].folderID, event.UserID)
	if err != nil {
		return err
	}
	if access != storage.Owner {
		p.sessions[event.UserID].status = statusOK
		return p.tg.SendMessage(event.ChatID, msgIncorrectAccessLvl)
	}

	_, err = p.storage.FolderID(ctx, event.UserID, event.Text)
	if err == nil {
		return p.tg.SendMessage(event.ChatID, msgCantRename)
	}
	if err != storage.ErrNoFolders {
		return err
	}

	err = p.storage.RenameFolder(ctx, p.sessions[event.UserID].folderID, event.Text)
	if err != nil {
		return err
	}

	return p.tg.SendMessage(event.ChatID, msgFolderRenamed)
}

func (p *Processor) deleteFolder(ctx context.Context, ChatID int, UserID int) error {

	access, err := p.storage.AccessLevelByUserID(ctx, p.sessions[UserID].folderID, UserID)
	if err != nil {
		return err
	}

	if access == storage.Owner {
		err = p.storage.RemoveFolder(ctx, p.sessions[UserID].folderID)
	} else {
		err = p.storage.DeleteAccess(ctx, UserID, p.sessions[UserID].folderID)
	}

	return p.tg.SendMessage(ChatID, msgFolderDeleted)
}
