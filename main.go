package main

import (
	"fmt"
	"github.com/sQUARys/TestTaskMailGaner/app/celerySender"
	"github.com/sQUARys/TestTaskMailGaner/app/client-mail/mailControllers"
	"github.com/sQUARys/TestTaskMailGaner/app/client-mail/mailRepositories"
	"github.com/sQUARys/TestTaskMailGaner/app/client-mail/mailRouters"
	"github.com/sQUARys/TestTaskMailGaner/app/client-mail/mailServices"
	"github.com/sQUARys/TestTaskMailGaner/app/senders"
	"log"
	"net/http"
	"time"
)

func main() {
	celerySenderEmails := celerySender.New()
	celerySenderEmails.Start()

	mailRepo := mailRepositories.New()

	mailService := mailServices.New(mailRepo)

	err := mailService.Start()
	if err != nil {
		log.Println(fmt.Errorf("Error : %w .\n", err))
	}
	emails, _ := mailService.GetEmails()
	fmt.Println("MAILSERVICE : ", emails)

	mailController := mailControllers.New(mailService)

	mailRouter := mailRouters.New(mailController)

	mailRouter.SetRoutes()

	mailServer := http.Server{
		ReadTimeout:  50 * time.Second,
		WriteTimeout: 50 * time.Second,
		Addr:         ":8081",
		Handler:      mailRouter.Router,
	}

	controller := senders.New()

	go controller.StartSending(mailRepo)

	err = mailServer.ListenAndServe()
	if err != nil {
		log.Println("Error in main : ", err)
		return
	}

}
