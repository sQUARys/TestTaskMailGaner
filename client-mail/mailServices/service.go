package mailServices

import (
	"github.com/sQUARys/TestTaskMailGaner/models"
	"sync"
)

type Service struct {
	Repo  usersRepository
	Cache emailsCache
	sync.RWMutex
}

type usersRepository interface {
	AddMail(mail models.Mail) error
	GetMails() ([]models.Mail, error)
	GetMailById(id int) (models.Mail, error)
	GetMailsByEmail(email string) ([]models.Mail, error)
}

type emailsCache interface {
	AddUserEmail(emailAddress string) error
	GetEmail(key string) (models.EmailAddress, error)
	GetEmails() ([]models.EmailAddress, error)
}

func New(repository usersRepository, cache emailsCache) *Service {
	serv := Service{
		Repo:  repository,
		Cache: cache,
	}
	return &serv
}

func (service *Service) GetMails() ([]models.Mail, error) {
	mails, err := service.Repo.GetMails()
	return mails, err
}

func (service *Service) AddMail(mail models.Mail) error {
	err := service.Repo.AddMail(mail)
	return err
}

func (service *Service) GetMailById(id int) (models.Mail, error) {
	mail, err := service.Repo.GetMailById(id)
	return mail, err
}

func (service *Service) GetMailsByEmail(email string) ([]models.Mail, error) {
	mails, err := service.Repo.GetMailsByEmail(email)
	return mails, err
}
