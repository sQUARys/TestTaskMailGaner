package main

import (
	"github.com/sQUARys/TestTaskMailGaner/app/controllers"
	"github.com/sQUARys/TestTaskMailGaner/client-mail/mailCache"
	"github.com/sQUARys/TestTaskMailGaner/client-mail/mailControllers"
	"github.com/sQUARys/TestTaskMailGaner/client-mail/mailRepositories"
	"github.com/sQUARys/TestTaskMailGaner/client-mail/mailRouters"
	"github.com/sQUARys/TestTaskMailGaner/client-mail/mailServices"
	"log"
	"net/http"
	"time"
)

func main() {
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

	controller := controllers.New()

	go controller.StartSending(emailCache)

	err := mailServer.ListenAndServe()
	if err != nil {
		log.Println("Error in main : ", err)
		return
	}

}
