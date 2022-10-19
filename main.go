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
	mailRepo := mailRepositories.New()
	mailService := mailServices.New(mailRepo)
	mailController := mailControllers.New(mailService)
	controller := senders.New()

	go func() {
		err := mailService.Start()
		if err != nil {
			log.Println(fmt.Errorf("Error : %w .\n", err))
		}
		celerySenderEmails.Start()
		controller.StartSending(mailRepo)
	}()

	mailRouter := mailRouters.New(mailController)
	mailRouter.SetRoutes()

	mailServer := http.Server{
		ReadTimeout:  50 * time.Second,
		WriteTimeout: 50 * time.Second,
		Addr:         ":8081",
		Handler:      mailRouter.Router,
	}

	err := mailServer.ListenAndServe()
	if err != nil {
		log.Println("Error in main : ", err)
		return
	}

}
