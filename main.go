package main

import (
	"github.com/sQUARys/TestTaskMailGaner/app/celerySender"
	"github.com/sQUARys/TestTaskMailGaner/app/client-mail/mailCache"
	"github.com/sQUARys/TestTaskMailGaner/app/client-mail/mailControllers"
	"github.com/sQUARys/TestTaskMailGaner/app/client-mail/mailRepositories"
	"github.com/sQUARys/TestTaskMailGaner/app/client-mail/mailRouters"
	"github.com/sQUARys/TestTaskMailGaner/app/client-mail/mailServices"
	"github.com/sQUARys/TestTaskMailGaner/app/models"
	"github.com/sQUARys/TestTaskMailGaner/app/senders"
	"log"
	"net/http"
	"time"
)

func main() {
	celerySender := celerySender.New()
	celeryMail := models.Mail{
		To:      "celeryEmail1@mail.ru",
		From:    "roman.kocenko@mail.ru",
		Message: "Hello from celery, bro)",
	}

	erro := celerySender.SendMessageWithTime(celeryMail, 5)
	if erro != nil {
		log.Println(erro)
	}

	emailCache := mailCache.New()
	emailCache.AddUserEmail("first@mail.ru")
	emailCache.AddUserEmail("second@mail.ru")
	emailCache.AddUserEmail("third@mail.ru")

	mailRepo := mailRepositories.New()

	mailService := mailServices.New(mailRepo, emailCache)

	mailController := mailControllers.New(mailService)

	mailRouter := mailRouters.New(mailController)

	mailRouter.SetRoutes()

	mailServer := http.Server{
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 3 * time.Second,
		Addr:         ":8081",
		Handler:      mailRouter.Router,
	}

	controller := senders.New()

	go controller.StartSending(emailCache)

	err := mailServer.ListenAndServe()
	if err != nil {
		log.Println("Error in main : ", err)
		return
	}

}
