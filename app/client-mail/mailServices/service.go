package mailServices

import (
	"github.com/sQUARys/TestTaskMailGaner/app/models"
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
	GetMailsByEmail(email string) ([]models.Mail, error)
	AddUserEmail(recipientEmailAddress string) error
	GetEmails() ([]models.EmailAddress, error)
}

func New(repository usersRepository) *Service {
	serv := Service{
		Repo: repository,
	}
	return &serv
}

func (service *Service) Start() error {
	for _, email := range models.UsersEmails {
		err := service.AddUserEmail(email)
		if err != nil {
			return err
		}
	}
	return nil
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

func (service *Service) AddUserEmail(recipientEmailAddress string) error {
	err := service.Repo.AddUserEmail(recipientEmailAddress)
	return err
}

func (service *Service) GetEmails() ([]models.EmailAddress, error) {
	emails, err := service.Repo.GetEmails()
	return emails, err
}
