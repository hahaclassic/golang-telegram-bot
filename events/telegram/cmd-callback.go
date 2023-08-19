package telegram

import "github.com/hahaclassic/golang-telegram-bot.git/lib/errhandling"

func (p *Processor) doCallbackCmd(text string, meta *CallbackMeta) (err error) {
	defer func() { err = errhandling.WrapIfErr("can't do callback cmd", err) }()

	err = p.tg.SendMessage(meta.ChatID, "callback epta")
	if err != nil {
		return err
	}

	// err = p.tg.AnswerCallbackQuery(meta.QueryID)
	// if err != nil {
	// 	return err
	// }

	return nil
}
