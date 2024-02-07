package telegram

import (
	"context"
	"errors"
	"log"
	"net/url"
	"strings"

	"github.com/hahaclassic/golang-telegram-bot.git/events"
	"github.com/hahaclassic/golang-telegram-bot.git/lib/errhandling"
	"github.com/hahaclassic/golang-telegram-bot.git/storage"
)

// text = text of the message
func (p *Processor) startCmd(event *events.Event) (err error) {

	defer func() {
		// В случае ошибки мы прерываем выполнение операции
		if err != nil {
			p.sessions[event.UserID].status = statusOK
		}
		// Отсутствие папок не является ошибкой, которую необходимо логировать.
		if err == ErrNoFolders {
			err = nil
		}
	}()

	event.Text = strings.TrimSpace(event.Text)
	log.Printf("got new command '%s' from '%d'", event.Text, event.UserID)

	if event.Text == CancelCmd {
		return p.tg.SendMessage(event.ChatID, msgNoCurrentOperation)
	}

	// Начало процесса добавление ссылки, если текст сообщения является ссылкой.
	if isAddCmd(event.Text) {
		p.sessions[event.UserID].url = event.Text
		p.sessions[event.UserID].currentOperation = GetNameCmd
		p.sessions[event.UserID].status = statusProcessing
		button := []string{"without a tag"}
		return p.tg.SendCallbackMessage(event.ChatID, msgEnterUrlName, button, button)
	}

	// Обработка однотактовых операций.
	switch event.Text {
	case StartCmd:
		return p.sendHello(event.ChatID)
	case RusHelpCmd:
		return p.sendRusHelp(event.ChatID)
	case HelpCmd:
		return p.sendHelp(event.ChatID)
	case RndCmd:
		return p.sendRandom(context.Background(), event.ChatID, event.UserID)
	}

	// Обработка сложных операций
	p.sessions[event.UserID].currentOperation = event.Text
	p.sessions[event.UserID].status = statusProcessing
	switch event.Text {
	case CreateFolderCmd:
		err = p.tg.SendMessage(event.ChatID, msgEnterFolderName)
	case ShowFolderCmd:
		err = p.chooseFolder(context.Background(), event.ChatID, event.UserID)
	case ChooseFolderForRenamingCmd:
		err = p.chooseFolder(context.Background(), event.ChatID, event.UserID)
	case DeleteFolderCmd:
		err = p.chooseFolder(context.Background(), event.ChatID, event.UserID)
	case ChooseLinkForDeletionCmd:
		err = p.chooseFolder(context.Background(), event.ChatID, event.UserID)
	case FeedbackCmd:
		err = p.tg.SendMessage(event.ChatID, msgEnterFeedback)
	// case ChangeTagCmd:
	// 	err = p.chooseFolder(context.Background(), chatID, userID)
	default:
		p.sessions[event.UserID].status = statusOK
		err = p.tg.SendMessage(event.ChatID, msgUnknownCommand)
	}

	return err
}

func (p *Processor) handleCmd(event *events.Event) (err error) {
	defer func() {
		// В случае ошибки мы прерываем выполнение операции
		if err != nil {
			p.sessions[event.UserID].status = statusOK
		}
		// Отсутствие папок не является ошибкой, которую необходимо логировать.
		if err == ErrNoFolders {
			err = nil
		}
	}()

	event.Text = strings.TrimSpace(event.Text)
	log.Printf("got new command '%s' from '%d'", event.Text, event.UserID)

	if event.Text == CancelCmd {
		return p.cancelOperation(event.ChatID, event.UserID)
	}

	if p.sessions[event.UserID].currentOperation != GetNameCmd {
		p.sessions[event.UserID].status = statusOK
	}
	switch p.sessions[event.UserID].currentOperation {
	case CreateFolderCmd:
		err = p.createFolder(context.Background(), event) // text == folderName
	case RenameFolderCmd:
		err = p.renameFolder(context.Background(), event) // text == folderName
	case FeedbackCmd:
		err = p.tg.SendMessage(event.ChatID, msgThanksForFeedback)
		err = p.logger.SendMessage(p.adminChatID, "#feedback\n\n"+event.Text)
	case GetNameCmd:
		if len(event.Text) > maxCallbackMsgLen {
			return p.tg.SendMessage(event.ChatID, msgLongMessage)
		}
		p.sessions[event.UserID].tag = event.Text
		p.sessions[event.UserID].currentOperation = SaveLinkCmd
		err = p.chooseFolder(context.Background(), event.ChatID, event.UserID)
	default:
		err = p.unknownCommandHelp(event.ChatID, event.UserID)
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
	case ShowFolderCmd:
		message += "Select the folder whose contents you want to see " + msgCancel
	case DeleteFolderCmd:
		message += "Select the folder you want to delete " + msgCancel
	default:
		message = msgUnexpectedCommand
	}

	return p.tg.SendMessage(chatID, message)
}

// event.Text == folderName
func (p *Processor) createFolder(ctx context.Context, event *events.Event) (err error) {
	defer func() { err = errhandling.WrapIfErr("can't create folder", err) }()

	_, err = p.storage.FolderID(ctx, event.UserID, event.Text)
	if err == nil {
		return p.tg.SendMessage(event.ChatID, msgFolderAlreadyExists)
	}

	var folder *storage.Folder

	i := 0
	for ; i < maxAttemts; i++ {
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

// Done
func (p *Processor) chooseFolder(ctx context.Context, chatID int, userID int) (err error) {
	defer func() {
		err = errhandling.WrapIfErr("can't do command: chooseFolder()", err)
	}()

	folders, err := p.storage.GetFolders(ctx, userID)
	if err != nil {
		return err
	}
	if len(folders[0]) == 0 {
		p.sessions[userID].status = statusOK
		return p.tg.SendMessage(chatID, msgNoFolders)
	}

	return p.tg.SendCallbackMessage(chatID, msgChooseFolder, folders[1], folders[1])
}

// event.Text == newFolderName
func (p *Processor) renameFolder(ctx context.Context, event *events.Event) (err error) {

	defer func() { err = errhandling.WrapIfErr("can't rename folder", err) }()

	folderID, err := p.storage.FolderID(ctx, event.UserID, p.sessions[event.UserID].folderName)
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

	_, err = p.storage.FolderID(ctx, event.UserID, event.Text)
	if err == nil {
		return p.tg.SendMessage(event.ChatID, msgCantRename)
	}
	if err != storage.ErrNoFolders {
		return err
	}

	err = p.storage.RenameFolder(ctx, folderID, event.Text)
	if err != nil {
		return errhandling.Wrap("can't rename folder", err)
	}

	return p.tg.SendMessage(event.ChatID, msgFolderRenamed)
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

	return p.tg.SendMessage(chatID, page)
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
