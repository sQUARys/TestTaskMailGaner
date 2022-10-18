package mailControllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/sQUARys/TestTaskMailGaner/client-mail/mailRepositories"
	"github.com/sQUARys/TestTaskMailGaner/models"
	"html"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"sync"
)

const (
	ok             = http.StatusOK
	serverInternal = http.StatusInternalServerError
	badRequest     = http.StatusBadRequest
	notFound       = http.StatusNotFound
)

type ErrorResponse struct {
	error string `json:"error"`
}

type MailController struct {
	Service mailService
	sync.RWMutex
}

type mailService interface {
	GetMails() ([]models.Mail, error)
	GetMailById(id int) (models.Mail, error)
}

func New(service mailService) *MailController {
	return &MailController{
		Service: service,
	}
}

func (ctr *MailController) MailHandler(w http.ResponseWriter, r *http.Request) {
	ctr.RLock()
	defer ctr.RUnlock()
	mailRepository := mailRepositories.New()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		ErrorHandler(w, err, serverInternal)
	}

	mailRepository.AddMail(models.Mail{
		From:    "from",
		To:      "to",
		IsRead:  false,
		Message: string(body),
	})

	w.WriteHeader(200)
	w.Write(body)

}

func (ctr *MailController) GetMailById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")

	vars := mux.Vars(r)
	idString := vars["message-id"]

	idInt, err := strconv.Atoi(idString)
	if err != nil {
		ErrorHandler(w, err, serverInternal)
		return
	}

	mail, err := ctr.Service.GetMailById(idInt)
	if err != nil {
		fmt.Println(err)
		ErrorHandler(w, err, serverInternal)
		return
	}

	//fmt.Fprintln(w, mail)
	fmt.Println(mail.IsRead)

}

func (ctr *MailController) GetMails(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	mails, err := ctr.Service.GetMails()

	if err != nil {
		ErrorHandler(w, err, serverInternal)
	}

	for _, mailHTML := range mails {
		tpl, err := template.New("card.html").ParseFiles("app/templates/card.html")
		if err != nil {
			return
		}

		buf := &bytes.Buffer{}
		tpl.Execute(buf, mailHTML)

		fmt.Fprintln(w, html.UnescapeString(buf.String()))
	}
}

func ErrorHandler(w http.ResponseWriter, err error, statusCode int) {
	errorResponse := ErrorResponse{
		error: err.Error(),
	}

	errStrJSON, errMarshal := json.Marshal(&errorResponse)
	if errMarshal != nil {
		log.Println(errMarshal)
		return
	}

	w.WriteHeader(statusCode)
	_, errorWriting := w.Write(errStrJSON)
	if errorWriting != nil {
		log.Println(errorWriting)
		return
	}
}

//func SendOkMessage(w http.ResponseWriter, action string) {
//	message := fmt.Sprintf("Action %s done succesful.", action)
//	responseString := users.ResponseOK{
//		Message: message,
//	}
//
//	responseJSON, err := json.Marshal(responseString)
//	if err != nil {
//		ErrorHandler(w, err, http.StatusInternalServerError)
//	}
//
//	_, err = w.Write(responseJSON)
//	if err != nil {
//		ErrorHandler(w, err, http.StatusInternalServerError)
//	}
//}
