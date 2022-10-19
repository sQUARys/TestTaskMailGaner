package mailCache

import (
	"encoding/json"
	"github.com/go-redis/redis"
	"github.com/sQUARys/TestTaskMailGaner/app/models"
	"log"
	"strconv"
	"sync"
)

type Cache struct {
	Client         *redis.Client
	SequenceNumber int
	*sync.RWMutex
}

func New() *Cache {
	client := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	})

	_, err := client.Ping().Result()
	if err != nil {
		log.Fatalln(err)
	}

	return &Cache{
		Client:         client,
		SequenceNumber: 0,
	}
}

func (c *Cache) AddUserEmail(recipientEmailAddress string) error {
	c.RLock()
	defer c.RUnlock()

	email := models.EmailAddress{
		Address: recipientEmailAddress,
	}

	c.SequenceNumber++

	emailJSON, err := json.Marshal(email)
	if err != nil {
		return err
	}

	key := strconv.Itoa(c.SequenceNumber)
	err = c.Client.SetNX(key, emailJSON, 0).Err()
	if err != nil {
		return err
	}
	return nil
}

func (c *Cache) GetEmail(key string) (models.EmailAddress, error) {
	transactionJSON, err := c.Client.Get(key).Result()
	if err != nil {
		return models.EmailAddress{}, err
	}
	var transaction models.EmailAddress
	err = json.Unmarshal([]byte(transactionJSON), &transaction)
	return transaction, err
}

func (c *Cache) GetEmails() ([]models.EmailAddress, error) {
	iter := c.Client.Scan(0, "", 0).Iterator()
	var emails []models.EmailAddress

	for iter.Next() {
		email, err := c.GetEmail(iter.Val())
		if err != nil {
			return nil, err
		}
		emails = append(emails, email)
	}

	if err := iter.Err(); err != nil {
		return nil, err
	}

	return emails, nil
}
