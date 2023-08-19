package telegram

const msgHelp = `I can save and keep you pages. 

just enter the link and I'll save it.

to get a random link, enter /rnd`

const msgHello = "Hi there!\n\n" + msgHelp

const (
	msgUnknownCommand      = "Unknown command 🤔"
	msgNoSavedPages        = "You have no saved pages 😢"
	msgSaved               = "Saved! 👌"
	msgAlreadyExists       = "You already have this page in your list 😌"
	msgFolderAlreadyExists = "This folder already exists 😌"
	msgFolderNotExists     = "This folder doesn't exist 🥺"

	msgEnterFolderName = "Enter the folder name"
	msgEnterLink       = "Enter the link"
)

const (
	HelpCmd  = "/help"
	StartCmd = "/start"

	DeleteLinkCmd = "/delete_link" // Удаляет ссылку из нужной папки
	SaveLink      = "/save"        // Сохраняет ссылку
	//ChangeFolderCmd = "/change"      // Меняет местонахождение ссылки
	RndCmd = "/rnd" // Скидывает случайную ссылку

	ShowFolderCmd   = "/folder"        // Показывает содержимое папки
	CreateFolderCmd = "/create"        // Создает новую папку
	DeleteFolderCmd = "/delete_folder" // Удаляет папку
	RenameFolderCmd = "/rename"        // Изменяет название папки
)
