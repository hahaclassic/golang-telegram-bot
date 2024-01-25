package telegram

// При статусе ОК сессия будет удалятся из карты. (Возможно, будет удаляется через некоторое время)
type Session struct {
	currentOperation string
	url              string
	name             string
	folder           string
	status           bool
}

func (p *Processor) changeSessionData(userID int, new *Session) {
	p.sessions[userID] = new
}

func (p *Processor) changeSessionURL(userID int, url string) {
	session := p.sessions[userID]
	session.url = url
	p.changeSessionData(userID, session)
}

func (p *Processor) changeSessionStatus(userID int, status bool) {
	session := p.sessions[userID]
	session.status = status
	p.changeSessionData(userID, session)
}

func (p *Processor) changeSessionName(userID int, name string) {
	session := p.sessions[userID]
	session.name = name
	p.changeSessionData(userID, session)
}

func (p *Processor) changeSessionOperation(userID int, operation string) {
	session := p.sessions[userID]
	session.currentOperation = operation
	p.changeSessionData(userID, session)
}

func (p *Processor) changeSessionFolder(userID int, folder string) {
	session := p.sessions[userID]
	session.folder = folder
	p.changeSessionData(userID, session)
}
