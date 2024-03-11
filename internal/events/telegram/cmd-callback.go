package telegram

import (
	"context"
	"strings"

	"github.com/hahaclassic/golang-telegram-bot.git/internal/events"
	"github.com/hahaclassic/golang-telegram-bot.git/internal/storage"
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

	if isCallbackOperation(meta) {
		return p.handleCallbackOperation(event, meta)
	}

	if p.sessions[event.UserID].currentOperation == DeleteLinkCmd {
		p.sessions[event.UserID].tag = meta.Data
	} else {
		p.sessions[event.UserID].folderID = meta.Data
	}

	switch p.sessions[event.UserID].currentOperation {

	case ChooseFolderForRenamingCmd:
		p.sessions[event.UserID].currentOperation = RenameFolderCmd
		// return p.tg.SendMessage(event.ChatID, msgEnterNewFolderName)
		return p.tg.EditMessage(event.ChatID, p.sessions[event.UserID].lastMessageID, msgEnterNewFolderName, nil)

	case ChooseLinkForDeletionCmd:
		p.sessions[event.UserID].currentOperation = DeleteLinkCmd
		return p.chooseLinkForDeletion(context.Background(), event.ChatID, event.UserID)

	case GetNameCmd:
		p.sessions[event.UserID].currentOperation = SaveLinkCmd
		p.sessions[event.UserID].tag = p.sessions[event.UserID].url
		err = p.chooseFolder(context.Background(), event.ChatID, event.UserID)

	case KeyCmd:
		return p.showKeys(context.Background(), event.ChatID, event.UserID)

	case SaveLinkCmd:
		p.sessions[event.UserID].currentOperation = DoneCmd
		return p.savePage(context.Background(), event.ChatID, event.UserID)

	case ShowFolderCmd:
		p.sessions[event.UserID].currentOperation = DoneCmd
		return p.showFolder(context.Background(), event.ChatID, event.UserID)

	case DeleteFolderCmd:
		p.sessions[event.UserID].currentOperation = DoneCmd
		return p.deleteFolder(context.Background(), event.ChatID, event.UserID)

	case DeleteLinkCmd:
		p.sessions[event.UserID].currentOperation = DoneCmd
		return p.deleteLink(context.Background(), event.ChatID, event.UserID)
	}

	return nil
}

func isCallbackOperation(meta *CallbackMeta) bool {
	strOperation := strings.Split(meta.Data, ",")[0]
	return ToOperation(strOperation) != UndefCmd
}

func (p *Processor) handleCallbackOperation(event *events.Event, meta *CallbackMeta) error {
	data := strings.Split(meta.Data, ",")
	operation := ToOperation(data[0])

	switch operation {
	case GetAccessCmd:
		return p.setAccess(context.Background(), event.ChatID, meta.Data, event.Text)

	case GoBackCmd:
		p.goBack(context.Background(), event)

	case ChooseForCreationKeyCmd:
		return p.chooseAccessLvl(event.ChatID, event.UserID, CreateKeyCmd)

	case ChooseForDeletionKeyCmd:
		return p.chooseAccessLvl(event.ChatID, event.UserID, DeleteKeyCmd)

	case CreateKeyCmd:
		return p.createKey(context.Background(), event.ChatID, event.UserID, storage.ToAccessLvl(data[1]))

	case DeleteKeyCmd:
		return p.deleteKey(context.Background(), event.ChatID, event.UserID, storage.ToAccessLvl(data[1]))
	}
	return nil
}

func (p *Processor) goBack(ctx context.Context, event *events.Event) (err error) {
	switch p.sessions[event.UserID].currentOperation {
	case CreateKeyCmd:
		p.sessions[event.UserID].currentOperation = KeyCmd
		err = p.showKeys(ctx, event.ChatID, event.UserID)
	case DeleteKeyCmd:
		p.sessions[event.UserID].currentOperation = KeyCmd
		err = p.showKeys(ctx, event.ChatID, event.UserID)
	case KeyCmd:
		p.sessions[event.UserID].currentOperation = KeyCmd
		err = p.chooseFolder(ctx, event.ChatID, event.UserID)
	}

	return err
}
