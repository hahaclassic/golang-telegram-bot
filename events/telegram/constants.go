package telegram

import "errors"

const msgHelp = `With this bot, you can store your important links and sort them by folders😊

To save the link:
1. Create a folder using /create
2. Enter the link (https://example.com)
3. Select the folder where you want to save the link

To view the contents of a folder:
1. Enter the /folder command
2. Select the desired folder

To delete a folder:
1. Enter the command /delete_folder
2. Select a folder
!!! BE CAREFUL !!! this command will delete the folder and all its contents without the possibility of recovery

Other commands:
/help - help about the bot
/rus_help - help in Russian
/rnd - output a random link from any folder

All commands are available in the menu next to the input field.
Productive work!

*** 
also under development:
/delete - delete links
/rename - rename folders
***`

const msgRusHelp = `С помощью данного бота ты можешь хранить свои важные ссылки и сортировать их по папкам😊

Чтобы сохранить ссылку:
1. Создайте папку с помощью /create
2. Введите ссылку (https://example.com)
3. Выберите папку, в которую хотите сохранить ссылку

Чтобы посмотреть содержимое папки:
1. Введите команду /folder
2. Выберите нужную папку

Чтобы удалить папку:
1. Введите команду /delete_folder
2. Выберите папку
!!! БУДЬТЕ ВНИМАТЕЛЬНЫ !!! данная команда удалит папку и все ее содержимое без возможности восстановления

Прочие команды:
/help - справка о боте
/rus_help - Справка на русском
/rnd - вывод случайной ссылки из любой папки

Все команды доступны в меню рядом с полем ввода.
Продуктивной работы!

*** также в процессе разработки:
/delete - удаление ссылок
/rename - переименование папок
***`

const msgHello = "Hi there!\n\n" + msgHelp

const (
	// Error
	msgUnknownCommand  = "Unknown command 🤔"
	msgFolderNotExists = "This folder doesn't exist 🥺"
	msgNoSavedPages    = "You have no saved pages 😢"
	msgNoFolders       = "No existing folders 😢"
	msgEmptyFolder     = "This folder is still empty 😢"

	// Warning
	msgFolderAlreadyExists = "This folder already exists 😌"
	msgAlreadyExists       = "You already have this page in your list 😌"

	// OK
	msgNewFolderCreated = "New Folder created 😇"
	msgSaved            = "Saved! 👌"
	msgFolderDeleted    = "Folder deleted 🫡"
	msgPageDeleted      = "Link deleted 🫡"

	// Input Suggestion
	msgChooseFolder    = "Choose folder"
	msgChooseLink      = "Choose link for deletion"
	msgEnterFolderName = "Enter the folder name"
)

// User commands
const (
	HelpCmd    = "/help"
	RusHelpCmd = "/rus_help"
	StartCmd   = "/start"

	ChooseLinkForDeletionCmd = "/delete" // Удаляет ссылку из нужной папки
	SaveLinkCmd              = "/save"   // Сохраняет ссылку 2
	//ChangeFolderCmd = "/change"      // Меняет местонахождение ссылки
	RndCmd = "/rnd" // Скидывает случайную ссылку

	ShowFolderCmd   = "/folder"        // Показывает содержимое папки 3
	CreateFolderCmd = "/create"        // Создает новую папку 1
	DeleteFolderCmd = "/delete_folder" // Удаляет папку
	RenameFolderCmd = "/rename"        // Изменяет название папки
)

// Internal commands
const (
	DeleteLinkCmd = "/delete_link"
)

var NoFoldersErr = errors.New("No existing folders")
