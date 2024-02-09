package telegram

const maxCallbackMsgLen = 60

const msgHelp = `With this bot, you can store your important links and sort them by folders😊

To save the link:
1. Create a folder using /create
2. Enter the link (https://example.com)
3. Select the folder where you want to save the link
(To save to an existing folder, just enter the link)

To view the contents of a folder:
1. Enter the show command
2. Select the desired folder

To delete a folder:
1. Enter the command /delete_folder
2. Select a folder
!!! BE CAREFUL !!! this command will delete the folder and all its contents without the possibility of recovery

To delete a link:
1. Enter the command /delete
2. Select a folder
3. Select a link
!!! BE CAREFUL !!! this command will delete the link without the possibility of recovery

To abort the operation, type /cancel

Other commands:
/help - help about the bot
/help_rus - help in Russian
/rename - rename folder
/rnd - output a random link from any folder

All commands are available in the menu next to the input field.
Productive work!`

const msgRusHelp = `С помощью данного бота ты можешь хранить свои важные ссылки и сортировать их по папкам😊

Чтобы сохранить ссылку:
1. Создайте папку с помощью /create
2. Введите ссылку (https://example.com)
3. Выберите папку, в которую хотите сохранить ссылку
(Чтобы сохранить в уже существующую папку, просто введите ссылку)

Чтобы посмотреть содержимое папки:
1. Введите команду /show
2. Выберите нужную папку

Чтобы удалить папку:
1. Введите команду /delete_folder
2. Выберите папку
!!! БУДЬТЕ ВНИМАТЕЛЬНЫ !!! данная команда удалит папку и все ее содержимое без возможности восстановления

Чтобы удалить ссылку:
1. Введите команду /delete
2. Выберите папку и все ее содержимое 
3. Выберите ссылку
!!! БУДЬТЕ ВНИМАТЕЛЬНЫ !!! данная команда удалит ссылку без возможности восстановления

Чтобы прервать операцию, введите /cancel

Прочие команды:
/help - справка о боте
/help_rus - Справка на русском
/rename - переименование папки
/rnd - вывод случайной ссылки из любой папки

Все команды доступны в меню рядом с полем ввода.
Продуктивной работы!`

const msgHello = "Hi there!\n\n" + msgHelp

const (
	// Error
	msgNoCurrentOperation = "At the moment, no operation is being performed 😶"
	msgUnknownCommand     = "Unknown command 🤔"
	msgUnexpectedCommand  = "Unexpected command 🤕"
	msgFolderNotExists    = "This folder doesn't exist 🥺"
	msgNoSavedPages       = "You have no saved pages 😢"
	msgNoFolders          = "No existing folders 😢"
	msgEmptyFolder        = "This folder is still empty 😢"
	msgCantRename         = "Cannot be renamed. A folder with this name already exists 😧"
	msgLongMessage        = "The message is too long, enter something shorter 🥴"
	msgIncorrectAccessLvl = "The operation has been stopped. Unfortunately, you don't have the right level of access 🔒"

	// Warning
	msgFolderAlreadyExists = "This folder already exists 😌"
	msgAlreadyExists       = "You already have this page in your list 😌"

	// OK
	msgNewFolderCreated   = "New Folder created 😇"
	msgSaved              = "Saved! 👌"
	msgFolderDeleted      = "Folder deleted 🫡"
	msgPageDeleted        = "Link deleted 🫡"
	msgFolderRenamed      = "Folder renamed 👌"
	msgOperationCancelled = "Operation cancelled 🤓"
	msgThanksForFeedback  = "Thank you for your help in improving our service! 🥺"

	// Input Suggestion
	msgChooseFolder       = "Choose folder"
	msgChooseLink         = "Choose link for deletion"
	msgEnterFolderName    = "Enter the folder name"
	msgEnterNewFolderName = "Enter new folder name"
	msgEnterUrlName       = "Enter short description (tag) for link"
	msgEnterFeedback      = "Write your feedback, ideas or suggestions. Don't worry, it's anonymous 💫"
)

const (
	HelpCmd     = "/help"
	RusHelpCmd  = "/help_rus"
	StartCmd    = "/start"
	CancelCmd   = "/cancel"
	FeedbackCmd = "/feedback"

	ChooseLinkForDeletionCmd = "/delete" // Удаляет ссылку из нужной папки
	//ChangeFolderCmd = "/move"      // Меняет местонахождение ссылки
	RndCmd = "/rnd" // Скидывает случайную ссылку
	// RenameLink = "/rename"
	// ChangeTagCmd = "/change_tag" // Изменение тега ссылки

	ShowFolderCmd              = "/show"          // Показывает содержимое папки
	CreateFolderCmd            = "/create"        // Создает новую папку
	DeleteFolderCmd            = "/delete_folder" // Удаляет папку
	ChooseFolderForRenamingCmd = "/rename"        // Изменяет название папки

	KeyCmd = "/key"
)

// Internal commands
const (
	SaveLinkCmd     = "/save"
	DeleteLinkCmd   = "/delete_link"
	RenameFolderCmd = "/rename_folder"
	GetNameCmd      = "/get_name"
	DeleteKeyCmd    = "/delete_key"
	CreateKeyCmd    = "/create_key"
)
