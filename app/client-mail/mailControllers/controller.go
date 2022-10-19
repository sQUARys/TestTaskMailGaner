package mailControllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/sQUARys/TestTaskMailGaner/app/models"
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
)

type MailController struct {
	Service mailService
	sync.RWMutex
}

type ErrorResponse struct {
	error string
}

type mailService interface {
	AddMail(mail models.Mail) error
	GetMails() ([]models.Mail, error)
	GetMailById(id int) (models.Mail, error)
	GetMailsByEmail(email string) ([]models.Mail, error)
}

func New(service mailService) *MailController {
	return &MailController{
		Service: service,
	}
}

func (ctr *MailController) MailHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	//ctr.RLock()
	//defer ctr.RUnlock()

	vars := mux.Vars(r)
	email := vars["email"]

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		ErrorHandler(w, err, serverInternal)
		return
	}

	err = ctr.Service.AddMail(models.Mail{
		To:      email,
		Message: string(body),
		IsRead:  false,
	})
	if err != nil {
		ErrorHandler(w, err, serverInternal)
		return
	}
}

func (ctr *MailController) GetMailsByEmail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	vars := mux.Vars(r)
	email := vars["email"]

	mails, err := ctr.Service.GetMailsByEmail(email)

	if err != nil {
		ErrorHandler(w, err, serverInternal)
		return
	}

	w.WriteHeader(ok)
	for _, mailHTML := range mails {
		err = WriteHTML(w, mailHTML, "card.html", "app/templates/card.html")
		if err != nil {
			log.Println(fmt.Sprintf("Error in  writing html. %w", err))
		}
	}

}

func (ctr *MailController) GetMailById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	vars := mux.Vars(r)
	idString := vars["message-id"]

	idInt, err := strconv.Atoi(idString)
	if err != nil {
		ErrorHandler(w, err, serverInternal)
		return
	}

	mail, err := ctr.Service.GetMailById(idInt)
	if err != nil {
		ErrorHandler(w, err, serverInternal)
		return
	}

	err = WriteHTML(w, mail, "description.html", "app/templates/description.html")
	if err != nil {
		log.Println(fmt.Sprintf("Error in  writing html. %w", err))
	}

}

func (ctr *MailController) GetMails(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	mails, err := ctr.Service.GetMails()

	if err != nil {
		ErrorHandler(w, err, serverInternal)
		return
	}

	w.WriteHeader(ok)
	for _, mailHTML := range mails {
		err = WriteHTML(w, mailHTML, "card.html", "app/templates/card.html")
		if err != nil {
			log.Println(fmt.Sprintf("Error in  writing html. %w", err))
		}
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

func WriteHTML(w http.ResponseWriter, mail interface{}, name string, path string) error {
	tpl, err := template.New(name).ParseFiles(path)
	if err != nil {
		return err
	}

	buf := &bytes.Buffer{}
	err = tpl.Execute(buf, mail)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintln(w, html.UnescapeString(buf.String()))
	if err != nil {
		return err
	}

	return nil
}
