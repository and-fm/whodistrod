package logging

type LogModel struct {
	AppName   string
	Env       string
	Method    string
	URI       string
	RemoteIP  string
	Status    int
	Latency   string
	Error     string
	SessionId string
	UserId    int
	UserEmail string
	Username  string
}
