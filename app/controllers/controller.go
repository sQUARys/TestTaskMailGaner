package controllers

import (
	"bytes"
	"fmt"
	"github.com/sQUARys/TestTaskMailGaner/models"
	"html/template"
	"net/http"
	"os"
	"time"
)

type Controller struct {
}

func New() *Controller {
	return &Controller{}
}

func (ctr *Controller) StartSending(messages []models.Message) {
	tpl, err := template.New("message.html").ParseFiles("app/templates/message.html")
	if err != nil {
		return
	}

	buf := &bytes.Buffer{}

	for _, message := range messages {
		fmt.Println("MESSAGE : ", message)
		tpl.Execute(buf, message)

		bodyReader := bytes.NewReader(buf.Bytes())

		req, err := http.NewRequest(http.MethodPost, "http://localhost:8081/mail/", bodyReader)
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
