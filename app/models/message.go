package models

type Message struct {
	EmailToSend string
	Message     string
}

type EmailAddress struct {
	Address string
}

type Mail struct {
	MessageId int
	From      string
	To        string
	Message   string
	IsRead    bool
}
