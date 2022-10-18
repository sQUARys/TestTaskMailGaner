package models

type Message struct {
	Name      string
	Age       int
	Message   string
	MessageId int
}

type Mail struct {
	MessageId int    `json:"message_id,omitempty"`
	From      string `json:"from,omitempty"`
	To        string `json:"to,omitempty"`
	Message   string `json:"message,omitempty"`
	IsRead    bool   `json:"isread,omitempty"`
}
