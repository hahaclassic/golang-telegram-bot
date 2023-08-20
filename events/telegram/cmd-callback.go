package telegram

import (
	"context"
	"log"
	"strings"

	conc "github.com/hahaclassic/golang-telegram-bot.git/lib/concatenation"
	"github.com/hahaclassic/golang-telegram-bot.git/lib/errhandling"
	"github.com/hahaclassic/golang-telegram-bot.git/storage"
)

func (p *Processor) doCallbackCmd(text string, meta *CallbackMeta) (err error) {
	defer func() {
		p.status = statusOK
		_ = p.tg.AnswerCallbackQuery(meta.QueryID)
		err = errhandling.WrapIfErr("can't do callback cmd", err)
	}()

	text = strings.TrimSpace(text)

	//err = p.tg.SendMessage(meta.ChatID, "callback epta")
	// if err != nil {
	// 	return err
	// }

	log.Println(text, meta.Message)

	switch p.currentOperation {
	case SaveLinkCmd:
		return p.savePage(context.Background(), meta, text)
	case ShowFolderCmd:
		return p.showFolder(context.Background(), meta, text)
	case DeleteFolderCmd:
		return p.deleteFolder(context.Background(), meta, text)
	}

	return nil
}

func (p *Processor) savePage(ctx context.Context, meta *CallbackMeta, folder string) (err error) {
	defer func() { err = errhandling.WrapIfErr("can't save page", err) }()

	page := &storage.Page{
		URL:    p.lastMessage,
		UserID: meta.UserID,
		Folder: folder,
	}

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

	result := conc.EnumeratedJoin(urls)
	if result == "" {
		return p.tg.SendMessage(meta.ChatID, msgEmptyFolder)
	}

	return p.tg.SendMessage(meta.ChatID, folder+":\n"+result)
}

func (p *Processor) deleteFolder(ctx context.Context, meta *CallbackMeta, folder string) error {

	err := p.storage.RemoveFolder(ctx, meta.UserID, folder)
	if err != nil {
		return errhandling.Wrap("can't delete folder", err)
	}

	return p.tg.SendMessage(meta.ChatID, msgFolderDeleted)
}
