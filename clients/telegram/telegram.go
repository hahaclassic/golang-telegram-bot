package tgClient

import (
	"bytes"
	"encoding/json"
	"errors"
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

var NoDataErr = errors.New("no data")

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

func (c *Client) SendCallbackMessage(chatID int, text string, list []string) error {
	buttons := [][]InlineKeyboardButton{}

	if list == nil || len(list) == 0 {
		return NoDataErr
	}

	for _, url := range list {
		inline := []InlineKeyboardButton{}
		inline = append(inline, InlineKeyboardButton{
			Text:         url,
			CallbackData: url,
		})
		buttons = append(buttons, inline)
	}

	replyMarkup := InlineKeyboardMarkup{
		InlineKeyboard: buttons}

	log.Println(replyMarkup) // LOG

	data := ReplyMessage{
		ChatID:      chatID,
		Text:        text,
		ReplyMarkup: replyMarkup,
	}

	// Get json
	EncodedData, err := json.Marshal(data)
	if err != nil {
		return errhandling.Wrap("can't get json", err)
	}

	_, err = c.doPostRequest(sendMessageMethod, EncodedData)
	if err != nil {
		return errhandling.Wrap("can't send a callback message", err)
	}

	return nil
}

func (c *Client) AnswerCallbackQuery(CallbackQueryID string) error {
	q := url.Values{}
	q.Add("callback_query_id", CallbackQueryID)

	_, err := c.doGetRequest(AnswerCallbackQueryMethod, q)
	if err != nil {
		return err
	}

	return nil
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
