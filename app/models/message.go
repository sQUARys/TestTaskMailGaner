package models

type Message struct {
	EmailToSend string
	Message     string
}

type EmailAddress struct {
	Address string
}

type Mail struct {
	MessageId int    `json:"message_id,omitempty"`
	From      string `json:"from,omitempty"`
	To        string `json:"to,omitemp2ty"`
	Message   string `json:"message,omitempty"`
	IsRead    bool   `json:"isread,omitempty"`
}
