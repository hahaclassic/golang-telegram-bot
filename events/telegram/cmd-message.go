package telegram

import (
	"context"
	"log"
	"net/url"
	"strings"

	"github.com/hahaclassic/golang-telegram-bot.git/events"
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
		_, err = p.tg.SendCallbackMessage(event.ChatID, msgEnterUrlName, button, button)
		return err
	}
	if isGetAccessCmd(event.Text) {
		p.sessions[event.UserID].status = statusOK
		return p.checkKey(context.Background(), event)
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
	case KeyCmd:
		err = p.chooseFolder(context.Background(), event.ChatID, event.UserID)
	case FeedbackCmd:
		err = p.tg.SendMessage(event.ChatID, msgEnterFeedback)
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

func isGetAccessCmd(text string) bool {
	return isKey(text)
}

func isKey(text string) bool {
	return len(text) > 3 && text[:3] == "KEY"
}
