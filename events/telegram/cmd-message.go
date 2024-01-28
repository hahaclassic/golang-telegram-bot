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

// func (p *Processor) doCmd(text string, chatID int, userID int) (err error) {

// 	defer func() {
// 		// В случае ошибки мы прерываем выполнение операции
// 		if err != nil {
// 			p.sessions[userID].status = statusOK
// 		}
// 		// Отсутствие папок не является ошибкой, которую необходимо логировать.
// 		if err == ErrNoFolders {
// 			err = nil
// 		}
// 	}()

// 	text = strings.TrimSpace(text)
// 	log.Printf("got new command '%s' from '%d'", text, userID)

// 	if text == CancelCmd {
// 		return p.cancelOperation(chatID, userID)
// 	}

// 	// Завершение/продолжение операции, если она в статусе обработки.
// 	if !p.sessions[userID].status {
// 		if p.sessions[userID].currentOperation != GetNameCmd {
// 			p.sessions[userID].status = statusOK
// 		}
// 		switch p.sessions[userID].currentOperation {
// 		case CreateFolderCmd:
// 			err = p.createFolder(context.Background(), chatID, userID, text) // text == folderName
// 		case RenameFolderCmd:
// 			err = p.renameFolder(context.Background(), chatID, userID, text) // text == folderName
// 		case GetNameCmd:
// 			if len(text) > maxCallbackMsgLen {
// 				return p.tg.SendMessage(chatID, msgLongMessage)
// 			}
// 			p.sessions[userID].name = text
// 			p.sessions[userID].currentOperation = SaveLinkCmd
// 			err = p.chooseFolder(context.Background(), chatID, userID)
// 		default:
// 			err = p.unknownCommandHelp(chatID, userID)
// 		}
// 		return err
// 	}

// 	// Добавление ссылки, если текст сообщения является ссылкой.
// 	if isAddCmd(text) {
// 		p.sessions[userID].url = text
// 		p.sessions[userID].currentOperation = GetNameCmd
// 		p.sessions[userID].status = statusProcessing
// 		return p.tg.SendCallbackMessage(chatID, msgEnterUrlName, []string{"without a tag"})
// 	}

// 	// // Начало выполнения новой операции.
// 	// // Обработка однотактовых операций.
// 	switch text {
// 	case StartCmd:
// 		return p.sendHello(chatID)
// 	case RusHelpCmd:
// 		return p.sendRusHelp(chatID)
// 	case HelpCmd:
// 		return p.sendHelp(chatID)
// 	case RndCmd:
// 		return p.sendRandom(context.Background(), chatID, userID)
// 	}

// 	// // Обработка сложных операций
// 	p.sessions[userID].currentOperation = text
// 	p.sessions[userID].status = statusProcessing
// 	switch text {
// 	case CreateFolderCmd:
// 		return p.tg.SendMessage(chatID, msgEnterFolderName)
// 	case ShowFolderCmd:
// 		return p.chooseFolder(context.Background(), chatID, userID)
// 	case ChooseFolderForRenamingCmd:
// 		return p.chooseFolder(context.Background(), chatID, userID)
// 	case DeleteFolderCmd:
// 		return p.chooseFolder(context.Background(), chatID, userID)
// 	case ChooseLinkForDeletionCmd:
// 		return p.chooseFolder(context.Background(), chatID, userID)
// 	default:
// 		p.sessions[userID].status = statusOK
// 		return p.tg.SendMessage(chatID, msgUnknownCommand)
// 	}
// }

// text = text of the message
func (p *Processor) startCmd(text string, chatID int, userID int) (err error) {

	defer func() {
		// В случае ошибки мы прерываем выполнение операции
		if err != nil {
			p.sessions[userID].status = statusOK
		}
		// Отсутствие папок не является ошибкой, которую необходимо логировать.
		if err == ErrNoFolders {
			err = nil
		}
	}()

	text = strings.TrimSpace(text)
	log.Printf("got new command '%s' from '%d'", text, userID)

	if text == CancelCmd {
		return p.tg.SendMessage(chatID, msgNoCurrentOperation)
	}

	// Начало процесса добавление ссылки, если текст сообщения является ссылкой.
	if isAddCmd(text) {
		p.sessions[userID].url = text
		p.sessions[userID].currentOperation = GetNameCmd
		p.sessions[userID].status = statusProcessing
		return p.tg.SendCallbackMessage(chatID, msgEnterUrlName, []string{"without a tag"})
	}

	// Обработка однотактовых операций.
	switch text {
	case StartCmd:
		return p.sendHello(chatID)
	case RusHelpCmd:
		return p.sendRusHelp(chatID)
	case HelpCmd:
		return p.sendHelp(chatID)
	case RndCmd:
		return p.sendRandom(context.Background(), chatID, userID)
	}

	// Обработка сложных операций
	p.sessions[userID].currentOperation = text
	p.sessions[userID].status = statusProcessing
	switch text {
	case CreateFolderCmd:
		err = p.tg.SendMessage(chatID, msgEnterFolderName)
	case ShowFolderCmd:
		err = p.chooseFolder(context.Background(), chatID, userID)
	case ChooseFolderForRenamingCmd:
		err = p.chooseFolder(context.Background(), chatID, userID)
	case DeleteFolderCmd:
		err = p.chooseFolder(context.Background(), chatID, userID)
	case ChooseLinkForDeletionCmd:
		err = p.chooseFolder(context.Background(), chatID, userID)
	default:
		p.sessions[userID].status = statusOK
		err = p.tg.SendMessage(chatID, msgUnknownCommand)
	}

	return err
}

func (p *Processor) handleCmd(text string, chatID int, userID int) (err error) {
	defer func() {
		// В случае ошибки мы прерываем выполнение операции
		if err != nil {
			p.sessions[userID].status = statusOK
		}
		// Отсутствие папок не является ошибкой, которую необходимо логировать.
		if err == ErrNoFolders {
			err = nil
		}
	}()

	text = strings.TrimSpace(text)
	log.Printf("got new command '%s' from '%d'", text, userID)

	if text == CancelCmd {
		return p.cancelOperation(chatID, userID)
	}

	if p.sessions[userID].currentOperation != GetNameCmd {
		p.sessions[userID].status = statusOK
	}
	switch p.sessions[userID].currentOperation {
	case CreateFolderCmd:
		err = p.createFolder(context.Background(), chatID, userID, text) // text == folderName
	case RenameFolderCmd:
		err = p.renameFolder(context.Background(), chatID, userID, text) // text == folderName
	case GetNameCmd:
		if len(text) > maxCallbackMsgLen {
			return p.tg.SendMessage(chatID, msgLongMessage)
		}
		p.sessions[userID].name = text
		p.sessions[userID].currentOperation = SaveLinkCmd
		err = p.chooseFolder(context.Background(), chatID, userID)
	default:
		err = p.unknownCommandHelp(chatID, userID)
	}

	return err
}

func (p *Processor) cancelOperation(chatID int, userID int) error {
	p.sessions[userID].status = statusOK
	return p.tg.SendMessage(chatID, msgOperationCancelled)
}

func (p *Processor) unknownCommandHelp(chatID int, userID int) error {

	var message string = msgUnexpectedCommand + "\n\n"
	var msgCancel string = "or enter /cancel to abort operation."

	switch p.sessions[userID].currentOperation {
	case ChooseFolderForRenamingCmd:
		message += "Select the folder you want to rename " + msgCancel
	case ChooseLinkForDeletionCmd:
		message += "Select the folder where you want to delete the link " + msgCancel
	case DeleteLinkCmd:
		message += "Select the link you want to delete "
	// case CreateFolderCmd:
	// 	message += "Enter new folder's name " + msgCancel
	case ShowFolderCmd:
		message += "Select the folder whose contents you want to see " + msgCancel
	case DeleteFolderCmd:
		message += "Select the folder you want to delete " + msgCancel
	default:
		message = msgUnexpectedCommand
	}

	return p.tg.SendMessage(chatID, message)
}

func (p *Processor) createFolder(ctx context.Context, chatID int, userID int, folder string) (err error) {
	defer func() { err = errhandling.WrapIfErr("can't create folder", err) }()

	ok, err := p.storage.IsFolderExist(ctx, userID, folder)
	if err == nil && ok {
		p.tg.SendMessage(chatID, msgFolderAlreadyExists)
	} else if err == nil {
		p.storage.NewFolder(ctx, userID, folder)
		p.tg.SendMessage(chatID, msgNewFolderCreated)
	}

	return err
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

func (p *Processor) renameFolder(ctx context.Context, chatID int, userID int, newFolder string) error {

	ok, err := p.storage.IsFolderExist(ctx, userID, newFolder)
	if err != nil {
		return errhandling.Wrap("can't rename folder", err)
	}
	if ok {
		return p.tg.SendMessage(chatID, msgCantRename)
	}

	err = p.storage.RenameFolder(ctx, userID, newFolder, p.sessions[userID].folder)
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
