package telegram

import (
	"context"
	"log"

	"github.com/hahaclassic/golang-telegram-bot.git/events"
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

	if len(meta.Data) > 7 && meta.Data[:7] == GetAccessCmd {
		return p.setAccess(context.Background(), event.ChatID, meta.Data, event.Text)
	}

	err = p.tg.DeleteMessage(event.ChatID, p.sessions[event.UserID].lastMessageID)
	if err != nil {
		return err
	}

	if meta.Data == GoBackCmd {
		return p.goBack(context.Background(), event)
	}

	if meta.Data == CreateKeyCmd || meta.Data == DeleteKeyCmd {
		p.sessions[event.UserID].currentOperation = meta.Data
	} else if storage.ToAccessLvl(meta.Data) != storage.Undefined {
		log.Println("YES")
		p.sessions[event.UserID].status = statusOK
		if p.sessions[event.UserID].currentOperation == CreateKeyCmd {
			return p.createKey(context.Background(), event.ChatID, event.UserID, storage.ToAccessLvl(meta.Data))
		} else {
			return p.deleteKey(context.Background(), event.ChatID, event.UserID, storage.ToAccessLvl(meta.Data))
		}
	} else if p.sessions[event.UserID].currentOperation == DeleteLinkCmd {
		p.sessions[event.UserID].tag = meta.Data
	} else {
		p.sessions[event.UserID].folderID = meta.Data
	}

	switch p.sessions[event.UserID].currentOperation {

	case ChooseFolderForRenamingCmd:
		p.sessions[event.UserID].currentOperation = RenameFolderCmd
		return p.tg.SendMessage(event.ChatID, msgEnterNewFolderName)

	case ChooseLinkForDeletionCmd:
		p.sessions[event.UserID].currentOperation = DeleteLinkCmd
		return p.chooseLinkForDeletion(context.Background(), event.ChatID, event.UserID)

	case GetNameCmd:
		p.sessions[event.UserID].currentOperation = SaveLinkCmd
		p.sessions[event.UserID].tag = p.sessions[event.UserID].url
		err = p.chooseFolder(context.Background(), event.ChatID, event.UserID)

	case KeyCmd:
		return p.showKeys(context.Background(), event.ChatID, event.UserID)

	case CreateKeyCmd:
		return p.chooseAccessLvl(event.ChatID, event.UserID)

	case DeleteKeyCmd:
		return p.chooseAccessLvl(event.ChatID, event.UserID)

	case SaveLinkCmd:
		p.sessions[event.UserID].status = statusOK
		return p.savePage(context.Background(), event.ChatID, event.UserID)

	case ShowFolderCmd:
		p.sessions[event.UserID].status = statusOK
		return p.showFolder(context.Background(), event.ChatID, event.UserID)

	case DeleteFolderCmd:
		p.sessions[event.UserID].status = statusOK
		return p.deleteFolder(context.Background(), event.ChatID, event.UserID)

	case DeleteLinkCmd:
		p.sessions[event.UserID].status = statusOK
		return p.deleteLink(context.Background(), event.ChatID, event.UserID)
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
