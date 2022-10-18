package controllers

import (
	"bytes"
	"fmt"
	"github.com/sQUARys/TestTaskMailGaner/client-mail/mailCache"
	"github.com/sQUARys/TestTaskMailGaner/models"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"
)

type Controller struct {
}

func New() *Controller {
	return &Controller{}
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

	for _, email := range emails {
		newMail := models.Mail{
			From:    "roman.kocenko@mail.ru",
			To:      email.Address,
			Message: "Push advertissment",
			IsRead:  false,
		}

		buf := &bytes.Buffer{}

		err = tpl.Execute(buf, newMail)
		if err != nil {
			log.Println(err)
			return
		}

		bodyReader := bytes.NewReader(buf.Bytes())

		req, err := http.NewRequest(http.MethodPost, "http://localhost:8081/mail/"+email.Address, bodyReader)
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

}
