package internal

type LogMessage struct {
	Text     string `json:"msg"`
	Level    string `json:"level"`
	Caller   string `json:"caller"`
	Time     string `json:"time"`
	Function string `json:"fName"`
}
