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
	AddMail(mail models.Mail) error
	GetMails() ([]models.Mail, error)
	GetMailById(id int) (models.Mail, error)
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

func (service *Service) GetMailById(id int) (models.Mail, error) {
	mail, err := service.Repo.GetMailById(id)
	return mail, err
}
