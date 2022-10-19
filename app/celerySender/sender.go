package celerySender

import (
	"bytes"
	"github.com/gocelery/gocelery"
	"github.com/gomodule/redigo/redis"
	"github.com/sQUARys/TestTaskMailGaner/app/models"
	"html/template"
	"time"
)

type Sender struct {
	TaskName string
	Client   *gocelery.CeleryClient
}

func New() *Sender {
	return &Sender{
		TaskName: "worker.sendMail",
		Client:   CreateCeleryClient(),
	}
}

func CreateCeleryClient() *gocelery.CeleryClient {
	redisPool := &redis.Pool{
		MaxIdle:     3,
		MaxActive:   0,
		IdleTimeout: 180 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.DialURL("redis://")
			if err != nil {
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}

	cli, _ := gocelery.NewCeleryClient(
		gocelery.NewRedisBroker(redisPool),
		&gocelery.RedisCeleryBackend{Pool: redisPool},
		1,
	)

	return cli
}

func (sender *Sender) SendMessageWithTime(mail models.Mail, sendAfterSeconds int) error {
	tpl, err := template.New("message.html").ParseFiles("app/templates/message.html")
	if err != nil {
		return err
	}

	buf := &bytes.Buffer{}

	err = tpl.Execute(buf, mail)
	if err != nil {
		return err
	}

	_, err = sender.Client.Delay(sender.TaskName, buf.String(), "http://localhost:8081/mail/"+mail.To, sendAfterSeconds)
	if err != nil {
		return err
	}
	return nil
}
