package main

import (
	"time"

	"github.com/gocelery/gocelery"
	"github.com/gomodule/redigo/redis"
)

// Run Celery Worker First!
// celery -A worker worker --loglevel=debug --without-heartbeat --without-mingle
func main() {

	// create redis connection pool
	redisPool := &redis.Pool{
		MaxIdle:     3,                 // maximum number of idle connections in the pool
		MaxActive:   0,                 // maximum number of connections allocated by the pool at a given time
		IdleTimeout: 240 * time.Second, // close connections after remaining idle for this duration
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

	// initialize celery client
	cli, _ := gocelery.NewCeleryClient(
		gocelery.NewRedisBroker(redisPool),
		&gocelery.RedisCeleryBackend{Pool: redisPool},
		1,
	)

	// prepare arguments
	//taskName := "worker.printHTML"
	////argA := rand.Intn(10)
	////argB := rand.Intn(10)
	//
	//// run task
	//asyncResult, err := cli.Delay(taskName, "HELLO")
	//if err != nil {
	//	panic(err)
	//}
	//
	//// get results from backend with timeout
	//res, err := asyncResult.Get(15 * time.Second)
	//if err != nil {
	//	panic(err)
	//}

	//log.Printf("result: %+v of type %+v", res, reflect.TypeOf(res))

	taskName := "worker.sendMail"

	_, err := cli.Delay(taskName, "<div style=\"display: flex ; flex-direction: column ; align-items: center ; justify-content: center\">\n    <h3>To : {{.To}}</h3> from your email address: {{.From}}\n    <h3>Send you a message:</h3><span>{{.Message}}</span>\n</div>", "http://localhost:8081/celery/email1@mail.ru", 5)
	if err != nil {
		panic(err)
	}

	// get results from backend with timeout

	//res, err := asyncResult.Get(10 * time.Second)
	//
	//result, err := cli.Delay("worker.getMail", res)
	//newResult, _ := result.Get(10 * time.Second)

	//fmt.Println("RES : ", res, err)
	//for res == nil {
	//	res, err = asyncResult.Get(10 * time.Second)
	//	fmt.Println("RES : ", res, err)
	//
	//}

	if err != nil {
		panic(err)
	}
	//log.Printf("result: %+v of type %+v", res, reflect.TypeOf(res))

}
