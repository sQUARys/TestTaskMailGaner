package main

import (
	"github.com/sQUARys/TestTaskMailGaner/app/controllers"
	"github.com/sQUARys/TestTaskMailGaner/client-mail/mailControllers"
	"github.com/sQUARys/TestTaskMailGaner/client-mail/mailRouters"
	"github.com/sQUARys/TestTaskMailGaner/models"
	"log"
	"net/http"
	"sync"
	"time"
)

func main() {
	mailController := mailControllers.New()
	mailRouter := mailRouters.New(mailController)

	mailRouter.SetRoutes()

	mailServer := http.Server{
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 3 * time.Second,
		Addr:         ":8081",
		Handler:      mailRouter.Router,
	}

	var wg *sync.WaitGroup

	controller := controllers.New()

	data := []models.Message{{Name: "Roman", Age: 18, Message: "Hello world."}, {Name: "Oleg", Age: 21, Message: "Creaste word."}}

	go controller.StartSending(data)

	err := mailServer.ListenAndServe()
	if err != nil {
		log.Println("Error in main : ", err)
		return
	}
	wg.Wait()

}
