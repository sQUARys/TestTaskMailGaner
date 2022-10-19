package models

const EmailFromWhichSendAllPushes = "roman.kocenko@mail.ru"

type PushNotification struct {
	To        string
	Push      string
	AfterTime int
}

var UsersEmails = []string{ // делаю для того, чтобы в бд юзерских имейлов было хоть что-либо, в реальном проекте эта бд уже будет дана
	"user1@gmail.com",
	"user2@gmail.com",
	"user3@mail.ru",
}

var Pushes = []string{
	"Hello, you won a bicycle!",
	"Whats'up, we can help you in education!",
	"Maybe you want to buy new shoes?",
}

var PushesForCelery = []PushNotification{
	{To: "celeryEmail1@mail.ru", Push: "Hello, you won a million dollars from Celery!", AfterTime: 5},
	{To: "celeryEmail2@mail.ru", Push: "Whats'up from celery, we can invite you to our university!", AfterTime: 10},
}
