package mailRouters

import (
	"github.com/gorilla/mux"
	"github.com/sQUARys/TestTaskMailGaner/client-mail/mailControllers"
	"log"
)

type Router struct {
	Router     *mux.Router
	Controller mailControllers.MailController
}

func New(controller *mailControllers.MailController) *Router {
	r := mux.NewRouter()
	return &Router{
		Controller: *controller,
		Router:     r,
	}
}

func (r *Router) SetRoutes() {
	log.Println("Mail routes start listening.")
	r.Router.HandleFunc("/mail/{email}", r.Controller.MailHandler).Methods("Post")
	r.Router.HandleFunc("/mail/sending-emails", r.Controller.GetMails).Methods("Get")
	r.Router.HandleFunc("/mail/got-emails/{email}", r.Controller.GetMailsByEmail).Methods("Get")
	r.Router.HandleFunc("/mail/message/{message-id:[0-9]+}", r.Controller.GetMailById).Methods("Get")
}

//:`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`
