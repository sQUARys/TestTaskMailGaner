package mailControllers

import (
	"encoding/json"
	"github.com/sQUARys/TestTaskMailGaner/client-mail/mailRepositories"
	"github.com/sQUARys/TestTaskMailGaner/models"
	"io/ioutil"
	"log"
	"net/http"
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
	sync.RWMutex
}

func New() *MailController {
	return &MailController{}
}

func (ctr *MailController) MailHandler(w http.ResponseWriter, r *http.Request) {
	mailRepository := mailRepositories.New()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		ErrorHandler(w, err, serverInternal)
	}
	//fmt.Println("Mail routers get : ", string(body))

	mailRepository.AddMessage(models.Mail{
		From:    "from",
		To:      "to",
		IsRead:  false,
		Message: string(body),
	})

	w.WriteHeader(200)
	w.Write(body)
	//
	//var usersTransfer users.TransferMoney
	//
	//err = json.Unmarshal(body, &usersTransfer)
	//if err != nil {
	//	ErrorHandler(w, err, serverInternal)
	//}
	//
	//if !ctr.Service.IsUserExisting(usersTransfer.IdOfSenderUser) || !ctr.Service.IsUserExisting(usersTransfer.IdOfRecipientUser) { // если в бд нет пользователя с таким id выводим ошибку
	//	ErrorHandler(w, err, notFound)
	//	return
	//}
	//
	//if usersTransfer.SendingAmount <= 0 {
	//	ErrorHandler(w, err, badRequest)
	//}
	//
	//ctr.RLock()
	//defer ctr.RUnlock() // для того чтобы при нескольких одновременно отправленных запросов на снятие, оно проходило последовательно и не было багов с балансом(чтобы не было гонки)
	//
	//err = ctr.Service.TransferMoney(usersTransfer)
	//if err != nil {
	//	ErrorHandler(w, err, serverInternal)
	//}
	//SendOkMessage(w, "transfer")

}

func (ctr *MailController) GetMails(w http.ResponseWriter, r *http.Request) {
	mailRepository := mailRepositories.New()
	mailRepository.GetMails()
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
