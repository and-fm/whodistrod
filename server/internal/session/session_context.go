package session

type SessionContext struct {
	SessionId string
	UserId    int
	Roles     []string
}
