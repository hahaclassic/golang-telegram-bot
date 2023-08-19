package tgClient

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
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
)

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

// func (c *Client) Updates(offset int, limit int) ([]Update, error) {
// 	q := url.Values{}
// 	q.Add("offset", strconv.Itoa(offset))
// 	q.Add("limit", strconv.Itoa(limit))

// 	data, err := c.doRequest(getUpdatesMethod, q)
// 	if err != nil {
// 		return nil, err
// 	}

// 	var res UpdatesResponse

// 	if err := json.Unmarshal(data, &res); err != nil {
// 		return nil, err
// 	}

// 	return res.Result, nil
// }

// func (c *Client) SendMessage(chatID int, text string) error {
// 	q := url.Values{}
// 	q.Add("chat_id", strconv.Itoa(chatID))
// 	q.Add("text", text)

// 	_, err := c.doRequest(sendMessageMethod, q)
// 	if err != nil {
// 		return errhandling.Wrap("can't send a message", err)
// 	}

// 	return nil
// }

// func (c *Client) doRequest(method string, query url.Values) (data []byte, err error) {

// 	defer func() { err = errhandling.WrapIfErr("can't do request", err) }()

// 	u := url.URL{
// 		Scheme: "https",
// 		Host:   c.host,
// 		Path:   path.Join(c.basePath, method),
// 	}

// 	req, err := http.NewRequest(http.MethodGet, u.String(), nil)

// 	if err != nil {
// 		return nil, err
// 	}

// 	req.URL.RawQuery = query.Encode()

// 	// sending a request to the telegram api
// 	resp, err := c.client.Do(req)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer func() { _ = resp.Body.Close() }()

// 	body, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return body, nil
// }

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

	data := StandardMessage{
		ChatID: chatID,
		Text:   text,
	}

	_, err := c.doPostRequest(sendMessageMethod, &data)
	if err != nil {
		return errhandling.Wrap("can't send a message", err)
	}

	return nil
}

// func (c *Client) SendCallbackMessage(chatID int, text string, list []string) error {
// 	buttons := [][]InlineKeyboardButton{}

// 	if list == nil || len(list) == 0 {
// 		return errors.New("no data")
// 	}

// 	for _, url := range list {
// 		inline := []InlineKeyboardButton{}
// 		inline = append(inline, InlineKeyboardButton{
// 			Text:         url,
// 			CallbackData: url,
// 		})
// 		buttons = append(buttons, inline)
// 	}

// 	replyMarkup := InlineKeyboardMarkup{
// 		InlineKeyboard: buttons}

// 	log.Println(replyMarkup) // LOG

// 	data := ReplyMessage{
// 		ChatID:      chatID,
// 		Text:        text,
// 		ReplyMarkup: replyMarkup,
// 	}

// 	_, err := c.doPostRequest(sendMessageMethod, &data)
// 	if err != nil {
// 		return errhandling.Wrap("can't send a message", err)
// 	}

// 	return nil
// }

func (c *Client) AnswerCallbackQuery(CallbackQueryID string) error {
	q := url.Values{}
	q.Add("callback_query_id", CallbackQueryID)

	_, err := c.doGetRequest(AnswerCallbackQueryMethod, q)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) doPostRequest(method string, query *StandardMessage) (data []byte, err error) {

	defer func() { err = errhandling.WrapIfErr("can't do request", err) }()

	log.Println(query.ChatID, query.Text)

	u := url.URL{
		Scheme: "https",
		Host:   c.host,
		Path:   path.Join(c.basePath, method),
	}

	EncodedData, err := json.Marshal(query)
	if err != nil {
		return nil, err
	}

	// Create new http post request
	req, err := http.NewRequest(http.MethodPost, u.String(), bytes.NewBuffer(EncodedData))

	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")

	// sending a request to the telegram api
	resp, err := c.client.Do(req)
	log.Println(resp.Status)
	if err != nil {
		return nil, err
	}

	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// var r UpdatesResponse

	// err = json.Unmarshal(body, &r)
	// log.Println(r.Result, r.Ok)

	// if err != nil {
	// 	return nil, err
	// }

	return body, nil
}

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
