package telegram

type Operation int

const (
	UndefCmd Operation = iota
	DoneCmd
	HelpCmd
	RusHelpCmd
	StartCmd
	CancelCmd
	FeedbackCmd

	ChooseLinkForDeletionCmd
	RndCmd

	ShowFolderCmd
	CreateFolderCmd
	DeleteFolderCmd
	ChooseFolderForRenamingCmd

	KeyCmd

	// Internal
	SaveLinkCmd
	DeleteLinkCmd

	RenameFolderCmd
	GetNameCmd

	ChooseForCreationKeyCmd
	ChooseForDeletionKeyCmd
	DeleteKeyCmd
	CreateKeyCmd
	GetAccessCmd

	GoBackCmd
)

func (op Operation) String() string {
	return []string{"/ok", "/help", "/help_rus", "/start", "/cancel", "/feedback",
		"/delete", "/rnd", "/show", "/create", "/delete_folder", "/rename", "key",
		"/save", "/delete_link", "/rename_folder", "/get_name", "/choose_for_creation_key",
		"/choose_for_deletion_key", "/delete_key", "/create_key", "/access", "/back"}[op]
}

func (op Operation) IsInternal() bool {
	return op >= SaveLinkCmd
}

func ToOperation(s string) Operation {
	for op := DoneCmd; op <= GoBackCmd; op++ {
		if s == op.String() {
			return op
		}
	}
	return DoneCmd
}
