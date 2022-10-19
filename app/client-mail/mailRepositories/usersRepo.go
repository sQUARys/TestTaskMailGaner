package mailRepositories

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/sQUARys/TestTaskMailGaner/app/models"
	"log"
	"sync"
	"time"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "myUser"
	password = "myPassword"
	dbname   = "myDb"

	connectionStringFormat = "host=%s port=%d user=%s password=%s dbname=%s sslmode=disable"

	//mail_table
	dbCreateMail       = `INSERT INTO mail_table("recipient","message", "isread") VALUES ($1, $2 , $3) RETURNING "message_id";`
	dbGetAllUsers      = "SELECT * FROM mail_table ORDER BY message_id"
	dbGetUserById      = "SELECT * FROM mail_table WHERE message_id = $1"
	dbGetMailByEmail   = "SELECT * FROM mail_table WHERE recipient = $1"
	dbUpdateStatusRead = "UPDATE mail_table SET isread=$1 WHERE message_id = $2"

	//users_email_table
	dbCreateUserEmail = `INSERT INTO users_email_table("email") VALUES ($1) RETURNING "email_id";`
	dbGetAllEmails    = "SELECT * FROM users_email_table ORDER BY email_id"
)

type Repository struct {
	DbStruct *sql.DB
	*sync.RWMutex
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

// Actions with mail_table

func (repo *Repository) AddMail(mail models.Mail) error {
	//repo.RLock()
	//defer repo.RUnlock()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := repo.DbStruct.ExecContext(
		ctx,
		dbCreateMail, mail.To, mail.Message, mail.IsRead,
	)

	if err != nil {
		return err
	}
	return nil
}

func (repo *Repository) GetMailById(id int) (mail models.Mail, err error) {
	//repo.RLock()
	//defer repo.RUnlock()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	tx, err := repo.DbStruct.BeginTx(ctx, nil)
	if err != nil {
		return
	}

	defer func() {
		if err != nil {
			err = tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	if err = tx.QueryRowContext(ctx, dbGetUserById, id).Scan(&mail.To, &mail.MessageId, &mail.Message, &mail.IsRead); err != nil {
		return
	}

	_, err = tx.ExecContext(ctx, dbUpdateStatusRead, true, mail.MessageId)
	if err != nil {
		return
	}

	return
}

func (repo *Repository) GetMailsByEmail(email string) ([]models.Mail, error) {
	rowsRs, err := repo.DbStruct.Query(dbGetMailByEmail, email)

	if err != nil {
		return []models.Mail{}, err
	}

	defer func() {
		err = rowsRs.Close()
		if err != nil {
			log.Println(err)
		}
	}()

	var mails []models.Mail

	for rowsRs.Next() {
		mail := models.Mail{}
		err = rowsRs.Scan(&mail.To, &mail.MessageId, &mail.Message, &mail.IsRead)
		if err != nil {
			return []models.Mail{}, err
		}
		mails = append(mails, mail)
	}

	if err = rowsRs.Err(); err != nil {
		return []models.Mail{}, err
	}

	return mails, err
}

func (repo *Repository) GetMails() ([]models.Mail, error) {
	rowsRs, err := repo.DbStruct.Query(dbGetAllUsers)

	if err != nil {
		return []models.Mail{}, err
	}

	defer func() {
		err = rowsRs.Close()
		if err != nil {
			log.Println(err)
		}
	}()

	var mails []models.Mail

	for rowsRs.Next() {
		mail := models.Mail{}
		err = rowsRs.Scan(&mail.To, &mail.MessageId, &mail.Message, &mail.IsRead)
		if err != nil {
			return []models.Mail{}, err
		}
		mails = append(mails, mail)
	}

	if err = rowsRs.Err(); err != nil {
		return []models.Mail{}, err
	}

	return mails, err
}

//Actions with users_email_table

func (repo *Repository) AddUserEmail(recipientEmailAddress string) error {
	//repo.RLock()
	//defer repo.RUnlock()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := repo.DbStruct.ExecContext(
		ctx,
		dbCreateUserEmail, recipientEmailAddress,
	)

	if err != nil {
		return err
	}

	return nil
}

func (repo *Repository) GetEmails() ([]models.EmailAddress, error) {
	rowsRs, err := repo.DbStruct.Query(dbGetAllEmails)

	if err != nil {
		return []models.EmailAddress{}, err
	}

	defer func() {
		err = rowsRs.Close()
		if err != nil {
			log.Println(err)
		}
	}()

	var emails []models.EmailAddress

	for rowsRs.Next() {
		var messageId int
		email := models.EmailAddress{}
		err = rowsRs.Scan(&messageId, &email.Address)
		if err != nil {
			return []models.EmailAddress{}, err
		}
		emails = append(emails, email)
	}

	if err = rowsRs.Err(); err != nil {
		return []models.EmailAddress{}, err
	}

	return emails, err
}
