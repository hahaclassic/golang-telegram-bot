package telegram

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/hahaclassic/golang-telegram-bot.git/storage"
)

var (
	decodeAccessDataErr = errors.New("Can't decode access data.")
)

type AccessData struct {
	FolderID    string
	FolderName  string
	AccessLevel storage.AccessLevel
	UserID      int
	Username    string
}

// returns folderID, userID, accessLvl from callbackData
// returns username, folderName from messageData
func decodeAccessData(callbackData string, message string) (*AccessData, error) {

	callbackParam := strings.Split(callbackData, ",")
	folderID := callbackParam[1]
	userID, err := strconv.Atoi(callbackParam[2])
	if err != nil {
		return nil, err
	}
	accessLevel := storage.ToAccessLvl(callbackParam[2])

	messageParam := strings.Split(message, "'")
	username, folderName := messageParam[1], messageParam[3]

	return &AccessData{
		FolderID:    folderID,
		FolderName:  folderName,
		AccessLevel: accessLevel,
		UserID:      userID,
		Username:    username,
	}, nil
}

func (data *AccessData) EncodeCallbackData() string {
	return fmt.Sprintf("%s,%s,%d,%s", GetAccessCmd, data.FolderID, data.UserID, data.AccessLevel)
}

func (data *AccessData) CreateMessage() string {
	return fmt.Sprintf("Предоставить ли доступ пользователю '%s' к папке '%s'?", data.Username, data.FolderName)
}
