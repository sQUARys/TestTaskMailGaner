package main

import (
	"log"
	"reflect"
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
	taskName := "worker.printHTML"
	//argA := rand.Intn(10)
	//argB := rand.Intn(10)

	// run task
	asyncResult, err := cli.Delay(taskName, "HELLO")
	if err != nil {
		panic(err)
	}

	// get results from backend with timeout
	res, err := asyncResult.Get(15 * time.Second)
	if err != nil {
		panic(err)
	}

	log.Printf("result: %+v of type %+v", res, reflect.TypeOf(res))

	taskName = "worker.add"

	asyncResult, err = cli.Delay(taskName, 4, 6)
	if err != nil {
		panic(err)
	}

	// get results from backend with timeout
	res, err = asyncResult.Get(5 * time.Second)
	if err != nil {
		panic(err)
	}
	log.Printf("result: %+v of type %+v", res, reflect.TypeOf(res))

}

//
//func main() {
//
//	// create redis connection pool
//	redisPool := &redis.Pool{
//		MaxIdle:     3,                 // maximum number of idle connections in the pool
//		MaxActive:   0,                 // maximum number of connections allocated by the pool at a given time
//		IdleTimeout: 240 * time.Second, // close connections after remaining idle for this duration
//		Dial: func() (redis.Conn, error) {
//			c, err := redis.DialURL("redis://")
//			if err != nil {
//				return nil, err
//			}
//			return c, err
//		},
//		TestOnBorrow: func(c redis.Conn, t time.Time) error {
//			_, err := c.Do("PING")
//			return err
//		},
//	}
//
//	// initialize celery client
//	cli, _ := gocelery.NewCeleryClient(
//		gocelery.NewRedisBroker(redisPool),
//		&gocelery.RedisCeleryBackend{Pool: redisPool},
//		1,
//	)
//
//	// prepare arguments
//	taskName := "worker.add"
//	//argA := rand.Intn(10)
//	//argB := rand.Intn(10)
//	messageHTML := "<div style=\"display: flex ; flex-direction: column ; align-items: center ; justify-content: center\">\n    <h3>To : {{.To}}</h3> from your email address: {{.From}}\n    <h3>Send you a message:</h3><span>{{.Message}}</span>\n</div>"
//	// run task
//	asyncResult, err := cli.Delay(taskName, messageHTML)
//	if err != nil {
//		panic(err)
//	}
//
//	// get results from backend with timeout
//	res, err := asyncResult.Get(30 * time.Second)
//
//	if err != nil {
//		panic(err)
//	}
//
//	log.Printf("result: %+v of type %+v", res, reflect.TypeOf(res))
//
//}

// Run Celery Worker First!
// celery -A worker worker --loglevel=debug --without-heartbeat --without-mingle
//func main() {
//	app := celery.NewApp()
//	app.Register(
//		"celery",
//		"mainQueue",
//		func(ctx context.Context, p *celery.TaskParam) error {
//			p.NameArgs("a", "b")
//			// Methods prefixed with Must panic if they can't find an argument name
//			// or can't cast it to the corresponding type.
//			// The panic doesn't affect other tasks execution; it's logged.
//			fmt.Println(p.MustInt("a") + p.MustInt("b"))
//			// Non-nil errors are logged.
//			return nil
//		},
//	)
//
//	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
//	defer cancel()
//
//	if err := app.Run(ctx); err != nil {
//		log.Printf("celery worker error: %v", err)
//	}
//}

//
// create redis connection pool
//redisPool := &redis.Pool{
//Dial: func() (redis.Conn, error) {
//c, err := redis.DialURL("redis://")
//if err != nil {
//return nil, err
//}
//return c, err
//},
//}
//
//// initialize celery client
//cli, _ := gocelery.NewCeleryClient(
//gocelery.NewRedisBroker(redisPool),
//&gocelery.RedisCeleryBackend{Pool: redisPool},
//1,
//)
//
//// prepare arguments
//taskName := "worker.add"
//argA := rand.Intn(10)
//argB := rand.Intn(10)
//
//// run task
//asyncResult, err := cli.Delay(taskName, argA, argB)
//if err != nil {
//panic(err)
//}
//
//// get results from backend with timeout
//res, err := asyncResult.Get(40 * time.Second)
//if err != nil {
//panic(err)
//}
//
//log.Printf("result: %+v of type %+v", res, reflect.TypeOf(res))
