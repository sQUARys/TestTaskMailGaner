package mailRepositories

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/sQUARys/TestTaskMailGaner/models"
	"log"
	"time"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "myUser"
	password = "myPassword"
	dbname   = "myDb"

	connectionStringFormat = "host=%s port=%d user=%s password=%s dbname=%s sslmode=disable"

	dbCreateMail        = `INSERT INTO mail_table("from" , "to", "message", "isread") VALUES ($1, $2 , $3 , $4) RETURNING "message_id";` //`INSERT INTO mail_table(from , to, message, isread) VALUES (%s , %s , %s , %t)` //RETURNING message_id
	dbGetAllUsers       = "SELECT * FROM mail_table"
	dbCreateUserRequest = `INSERT INTO "user_table"( "id") VALUES (%d)`
	dbUsersByIdRequest  = "SELECT * FROM user_table WHERE id = $1"
	dbUpdateJSON        = "UPDATE user_table SET balance=%2f WHERE id=%d"
)

type Repository struct {
	DbStruct *sql.DB
}

func New() *Repository {
	connectionString := fmt.Sprintf(connectionStringFormat, host, port, user, password, dbname)

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatalln(err)
	}

	repo := Repository{
		DbStruct: db,
	}

	err = db.Ping()
	if err != nil {
		log.Fatalln(err)
	}

	return &repo
}

func (repo *Repository) AddMessage(mail models.Mail) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	//если понадобится id сообщения
	//message_id := -1
	//
	//err := repo.DbStruct.QueryRow(dbCreateMail, mail.From, mail.To, mail.Message, mail.IsRead).Scan(&message_id)
	//if err != nil {
	//	panic(err)
	//}

	_, err := repo.DbStruct.ExecContext(
		ctx,
		dbCreateMail, mail.From, mail.To, mail.Message, mail.IsRead,
	)

	if err != nil {
		return err
	}
	return nil
}

func (repo *Repository) GetMails() {
	rowsRs, err := repo.DbStruct.Query(dbGetAllUsers)

	if err != nil {
		return
	}
	defer rowsRs.Close()

	var mails []models.Mail

	for rowsRs.Next() {
		mail := models.Mail{}
		err = rowsRs.Scan(&mail.From, &mail.To, &mail.MessageId, &mail.Message, &mail.IsRead)
		if err != nil {
			log.Println(err)
			return
		}
		mails = append(mails, mail)
	}

	fmt.Println("MAILS : ", mails)
	if err = rowsRs.Err(); err != nil {
		return
	}
}
