package senders

import (
	"bytes"
	"fmt"
	"github.com/sQUARys/TestTaskMailGaner/app/client-mail/mailCache"
	"github.com/sQUARys/TestTaskMailGaner/app/models"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"
)

type Controller struct {
	Pushes    []string
	FromEmail string
}

func New() *Controller {
	return &Controller{
		Pushes:    []string{"Hello, you won a bicycle!", "Whats'up, we can help you in education!"},
		FromEmail: "roman.kocenko@mail.ru",
	}
}

func (ctr *Controller) StartSending(cache *mailCache.Cache) {
	tpl, err := template.New("message.html").ParseFiles("app/templates/message.html")
	if err != nil {
		log.Println(err)
		return
	}

	emails, err := cache.GetEmails()
	if err != nil {
		return
	}

	for _, push := range ctr.Pushes {
		for _, email := range emails {
			newMail := models.Mail{
				From:    ctr.FromEmail,
				To:      email.Address,
				Message: push,
				IsRead:  false,
			}

			buf := &bytes.Buffer{}

			err = tpl.Execute(buf, newMail)
			if err != nil {
				log.Println(err)
				return
			}

			PostRequest(buf, email.Address)
		}
	}
}

func PostRequest(buf *bytes.Buffer, emailAddress string) {
	bodyReader := bytes.NewReader(buf.Bytes())

	req, err := http.NewRequest(http.MethodPost, "http://localhost:8081/mail/"+emailAddress, bodyReader)
	if err != nil {
		fmt.Printf("client: could not create request: %s\n", err)
		os.Exit(1)
	}
	req.Header.Set("Content-Type", "application/json")

	client := http.Client{
		Timeout: 30 * time.Second,
	}

	_, err = client.Do(req)
	if err != nil {
		fmt.Printf("client: error making http request: %s\n", err)
		os.Exit(1)
	}
}
