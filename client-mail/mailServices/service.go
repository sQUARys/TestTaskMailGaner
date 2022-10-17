package mailServices

import (
	"github.com/sQUARys/TestTaskMailGaner/models"
	"sync"
)

type Service struct {
	Repo usersRepository
	sync.RWMutex
}

type usersRepository interface {
	AddMessage(mail models.Mail) error
	GetMails() ([]models.Mail, error)
}

func New(repository usersRepository) *Service {
	serv := Service{
		Repo: repository,
	}
	return &serv
}

func (service *Service) GetMails() ([]models.Mail, error) {
	mails, err := service.Repo.GetMails()
	return mails, err
}
