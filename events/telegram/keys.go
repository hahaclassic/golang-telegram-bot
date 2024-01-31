package telegram

// Уровни доступа
// 0 - Owner: Может делать все, что и остальные + rename folder + delete folder
// 1 - write: add/delete links (always with confirmation)
// 2 - read with confirmation
// 3 - read

type AccessLevel int

const (
	Owner AccessLevel = iota
	Editor
	ConfirmedReader
	Reader
)

// Создание ключа для выбранной папки
// func (p *Processor) CreateKey(userID, folderID, access int) (string, error) {
// 	var key strings.Builder
// 	key.WriteString("KEY")
// 	key.WriteString(strconv.Itoa(userID) + "A")
// 	key.WriteString()

// }

// func (p *Processor) DeleteKey(userID, folderID int) error {}

// func (p *Processor) Check(key string) bool {}
