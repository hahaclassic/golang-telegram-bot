package tgClient

type UpdatesResponse struct {
	Ok     bool     `json:"ok"`
	Result []Update `json:"result"`
}

type PostRequestResponse struct {
	Ok     bool            `json:"ok"`
	Result returnedMessage `json:"result"`
}

// В запросе возвращаются поля 'message_id', 'from', 'chat', 'date', 'text' и т.д.,
// Но для нас необходим только message_id
type returnedMessage struct {
	MessageID int `json:"message_id"`
}

type Update struct {
	ID            int              `json:"update_id"`
	Message       *IncomingMessage `json:"message"`
	CallbackQuery *CallbackQuery   `json:"callback_query"`
}

type CallbackQuery struct {
	QueryID string           `json:"id"`
	From    From             `json:"from"`
	Message *IncomingMessage `json:"message"`
	Data    string           `json:"data"`
}

type IncomingMessage struct {
	Text string `json:"text"`
	From From   `json:"from"`
	Chat Chat   `json:"chat"`
}

type From struct {
	UserID   int    `json:"id"`
	Username string `json:"username"`
}

type Chat struct {
	ID int `json:"id"`
}

type ReplyMessage struct {
	ChatID      int                  `json:"chat_id"`
	Text        string               `json:"text"`
	ParseMode   string               `json:"parse_mode"`
	ReplyMarkup InlineKeyboardMarkup `json:"reply_markup,omitempty"`
}

type StandardMessage struct {
	ChatID int    `json:"chat_id"`
	Text   string `json:"text"`
	//ParseMode string `json:"parse_mode"`
}

type InlineKeyboardMarkup struct {
	InlineKeyboard [][]InlineKeyboardButton `json:"inline_keyboard"`
}

type InlineKeyboardButton struct {
	Text         string `json:"text"`
	CallbackData string `json:"callback_data"`
}
