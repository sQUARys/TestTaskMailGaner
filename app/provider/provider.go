package main

import (
	"github.com/gocelery/gocelery"
	"github.com/gomodule/redigo/redis"
	"time"
)

type exampleAddTask struct {
	a       int
	b       int
	message string
}

func (a *exampleAddTask) RunTask() (interface{}, error) {
	result := a.a + a.b
	return result, nil
}

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
		5, // number of workers
	)

	// register task
	cli.Register("worker.add", &exampleAddTask{})
	cli.Register("worker.send_mail", &exampleAddTask{})
	// start workers (non-blocking call)
	cli.StartWorker()

	// wait for client request
	time.Sleep(10 * time.Second)

	// stop workers gracefully (blocking call)
	cli.StopWorker()
}

//type exampleAddTask struct {
//	message string
//}
////
////func (a *exampleAddTask) RunTask() (interface{}, error) {
////	result := a.message
////	return result, nil
////}
////
////func main() {
////	// create redis connection pool
////	redisPool := &redis.Pool{
////		Dial: func() (redis.Conn, error) {
////			c, err := redis.DialURL("redis://")
////			if err != nil {
////				return nil, err
////			}
////			return c, err
////		},
////	}
////
////	// initialize celery client
////	cli, _ := gocelery.NewCeleryClient(
////		gocelery.NewRedisBroker(redisPool),
////		&gocelery.RedisCeleryBackend{Pool: redisPool},
////		5, // number of workers
////	)
////
////	//printHTML := func(message string) string {
////	//	return message
////	//}
////
////	// register task
////	cli.Register("worker.add", &exampleAddTask{})
////
////	// start workers (non-blocking call)
////	cli.StartWorker()
////
////	// wait for client request
////	time.Sleep(10 * time.Second)
////
////	// stop workers gracefully (blocking call)
////	cli.StopWorker()
////
////}

//func main() {
//	app := celery.NewApp()
//
//	err := app.Delay(
//		"celery",
//		"mainQueue",
//		"a",
//		3,
//	)
//	if err != nil {
//		log.Printf("failed to send mytask: %v", err)
//	}
//}

//
//// create redis connection pool
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
//// task
//add := func(a, b int) int {
//	return a + b
//}
//
//// register task
//cli.Register("worker", add)
//
//// context with cancelFunc to handle exit gracefully
//ctx, cancel := context.WithCancel(context.Background())
//
//// start workers (non-blocking call)
//cli.StartWorkerWithContext(ctx)
//
//// wait for client request
//time.Sleep(10 * time.Second)
//
//// stop workers by cancelling context
//cancel()
//
//// optional: wait for all workers to terminate
//cli.WaitForStopWorker()
