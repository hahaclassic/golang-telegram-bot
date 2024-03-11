package tgclient

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"

	"github.com/hahaclassic/golang-telegram-bot.git/lib/errhandling"
)

type Client struct {
	host     string
	basePath string
	client   http.Client
}

const (
	getUpdatesMethod          = "getUpdates"
	sendMessageMethod         = "sendMessage"
	AnswerCallbackQueryMethod = "answerCallbackQuery"
	deleteMessageMethod       = "deleteMessage"
	editMessageMethod         = "editMessageText"
)

var ErrNoData = errors.New("no data")
var ErrWrongData = errors.New("wrong data")

func New(host string, token string) *Client {
	return &Client{
		host:     host,
		basePath: newBasePath(token),
		client:   http.Client{},
	}
}

func newBasePath(token string) string {
	return "bot" + token
}

func (c *Client) Updates(offset int, limit int) ([]Update, error) {
	q := url.Values{}
	q.Add("offset", strconv.Itoa(offset))
	q.Add("limit", strconv.Itoa(limit))

	data, err := c.doGetRequest(getUpdatesMethod, q)
	if err != nil {
		return nil, err
	}

	var res UpdatesResponse

	if err := json.Unmarshal(data, &res); err != nil {
		return nil, err
	}

	return res.Result, nil
}

func (c *Client) SendMessage(chatID int, text string) error {

	data := OutputMessage{
		ChatID:    chatID,
		Text:      text,
		ParseMode: "HTML",
	}

	// Get json
	EncodedData, err := json.Marshal(data)
	if err != nil {
		return errhandling.Wrap("can't get json", err)
	}

	_, err = c.doPostRequest(sendMessageMethod, EncodedData)
	if err != nil {
		return errhandling.Wrap("can't send a message", err)
	}

	return nil
}

func CreateInlineKeyboardMarkup(buttonsText []string, callbackData []string) (*InlineKeyboardMarkup, error) {
	if len(buttonsText) == 0 || len(callbackData) == 0 {
		return nil, ErrNoData
	}
	if len(buttonsText) != len(callbackData) {
		return nil, ErrWrongData
	}

	buttons := [][]InlineKeyboardButton{}

	for i := 0; i < len(buttonsText); i++ {
		inline := []InlineKeyboardButton{}
		inline = append(inline, InlineKeyboardButton{
			Text:         buttonsText[i],
			CallbackData: callbackData[i],
		})
		buttons = append(buttons, inline)
	}

	replyMarkup := &InlineKeyboardMarkup{
		InlineKeyboard: buttons}

	return replyMarkup, nil
}

// SendCallbackMessage() returns messageID of the sent message and error
func (c *Client) SendCallbackMessage(chatID int, text string, replyMarkup *InlineKeyboardMarkup) (messageID int, err error) {

	data := OutputMessage{
		ChatID:      chatID,
		Text:        text,
		ParseMode:   "HTML",
		ReplyMarkup: *replyMarkup,
	}

	// Get json
	EncodedData, err := json.Marshal(data)
	if err != nil {
		return 0, errhandling.Wrap("can't get json", err)
	}

	bodyData, err := c.doPostRequest(sendMessageMethod, EncodedData)
	if err != nil {
		return 0, errhandling.Wrap("can't send a callback message", err)
	}

	var res PostRequestResponse
	if err := json.Unmarshal(bodyData, &res); err != nil {
		return 0, errhandling.Wrap("can't unmarshal json", err)
	}

	return res.Result.MessageID, nil
}

func (c *Client) EditMessage(chatID int, messageID int, text string, replyMarkup *InlineKeyboardMarkup) error {
	var reply InlineKeyboardMarkup
	if replyMarkup != nil {
		reply = *replyMarkup
	}

	data := OutputMessage{
		ChatID:      chatID,
		Text:        text,
		ParseMode:   "HTML",
		ReplyMarkup: reply,
	}

	EncodedData, err := json.Marshal(data)
	if err != nil {
		return errhandling.Wrap("can't get json", err)
	}

	bodyData, err := c.doPostRequest(editMessageMethod, EncodedData)
	if err != nil {
		return errhandling.Wrap("can't send a callback message", err)
	}

	var res PostRequestResponse
	if err := json.Unmarshal(bodyData, &res); err != nil {
		return errhandling.Wrap("can't unmarshal json", err)
	}

	return nil
}

func (c *Client) DeleteMessage(chatID int, messageID int) error {
	q := url.Values{}
	q.Add("chat_id", strconv.Itoa(chatID))
	q.Add("message_id", strconv.Itoa(messageID))

	_, err := c.doGetRequest(deleteMessageMethod, q)

	return err
}

func (c *Client) AnswerCallbackQuery(CallbackQueryID string) error {
	q := url.Values{}
	q.Add("callback_query_id", CallbackQueryID)

	_, err := c.doGetRequest(AnswerCallbackQueryMethod, q)

	return err
}

// doPostRequest() sends a post request to the server. Accepts data in json format
func (c *Client) doPostRequest(method string, jsonData []byte) (data []byte, err error) {

	defer func() { err = errhandling.WrapIfErr("can't do request", err) }()

	u := url.URL{
		Scheme: "https",
		Host:   c.host,
		Path:   path.Join(c.basePath, method),
	}

	// Create new http post request
	req, err := http.NewRequest(http.MethodPost, u.String(), bytes.NewBuffer(jsonData))

	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")

	// sending a request to the telegram api
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

// doPostRequest() sends a get request to the server. Accepts data in url.Values format
func (c *Client) doGetRequest(method string, query url.Values) (data []byte, err error) {

	defer func() { err = errhandling.WrapIfErr("can't do request", err) }()

	u := url.URL{
		Scheme: "https",
		Host:   c.host,
		Path:   path.Join(c.basePath, method),
	}

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)

	if err != nil {
		return nil, err
	}

	req.URL.RawQuery = query.Encode()

	// sending a request to the telegram api
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	data, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return data, nil
}
