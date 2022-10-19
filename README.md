# Тестовое задание на позицию стажера-бекендера
Микросервис для рассылки email сообщений на разные адреса.

## Задача

Написать на Go небольшой сервис отправки имейл рассылок.
Возможности сервиса:
 1. Отправка рассылок с использованием html макета и списка подписчиков.
 2. Отправка отложенных рассылок.
 3. Использование переменных в макете рассылки. (Пример: имейл получателя, имейл отправителя, текст сообщения)
 4. Отслеживание открытий писем.


## Стек используемых в сервисе технологий

* Golang
* PostgreSQL(две таблицы: для хранения информации о всех сообщениях отправленных; для хранения email адресов пользователей, которым должна прийти рассылка)
* Celery(для отложенной отправки писем)


## Начальные настройки

Перед тем, как запускать файл main.go, который запустит весь микросервис, необходимо создать таблицы в PostgreSQL

```SQL

CREATE TABLE mail_table(
    recipient TEXT not null,
    message_id INT  GENERATED ALWAYS AS IDENTITY,
    message TEXT not null,
    isRead bool DEFAULT false
);

CREATE TABLE users_email_table(
    email_id INT  GENERATED ALWAYS AS IDENTITY,
    email TEXT NOT NULL
);
```

Также, для работы с Celery нам необходимо его запустить. Переходим в папку [celerySender](https://github.com/sQUARys/TestTaskMailGaner/tree/master/app/celerySender)
Запускаем скрипт воркера
```
  python3 worker.py      
```
После этого, запускаем воркер Celery для взаимодействия с ним
```
  celery -A worker worker --loglevel=debug --without-heartbeat --without-mingle    
```

Теперь переходим в папку проекта и запускаем main.go файл. После его запуска, в Celery и users_email_table будут записаны
заданные по дефолту начальные условия, чтобы из данных сервисов можно было отправлять данные на сервер.
Данные начальные условия хранятся в [initialConditions](https://github.com/sQUARys/TestTaskMailGaner/blob/master/app/models/initialsConditions.go)


## Example

[GoCelery GoDoc](https://godoc.org/github.com/gocelery/gocelery) has good examples.<br/>
Also take a look at `example` directory for sample python code.

### GoCelery Worker Example

Run Celery Worker implemented in Go

```go
// create redis connection pool
redisPool := &redis.Pool{
  Dial: func() (redis.Conn, error) {
		c, err := redis.DialURL("redis://")
		if err != nil {
			return nil, err
		}
		return c, err
	},
}
// initialize celery client
cli, _ := gocelery.NewCeleryClient(
	gocelery.NewRedisBroker(redisPool),
	&gocelery.RedisCeleryBackend{Pool: redisPool},
	5, // number of workers
)
// task
add := func(a, b int) int {
	return a + b
}
// register task
cli.Register("worker.add", add)
// start workers (non-blocking call)
cli.StartWorker()
// wait for client request
time.Sleep(10 * time.Second)
// stop workers gracefully (blocking call)
cli.StopWorker()
```

### Python Client Example

Submit Task from Python Client

```python
from celery import Celery
app = Celery('tasks',
    broker='redis://localhost:6379',
    backend='redis://localhost:6379'
)
@app.task
def add(x, y):
    return x + y
if __name__ == '__main__':
    ar = add.apply_async((5456, 2878), serializer='json')
    print(ar.get())
```

### Python Worker Example

Run Celery Worker implemented in Python

```python
from celery import Celery
app = Celery('tasks',
    broker='redis://localhost:6379',
    backend='redis://localhost:6379'
)
@app.task
def add(x, y):
    return x + y
```

```bash
celery -A worker worker --loglevel=debug --without-heartbeat --without-mingle
```

### GoCelery Client Example

Submit Task from Go Client

```go
// create redis connection pool
redisPool := &redis.Pool{
  Dial: func() (redis.Conn, error) {
		c, err := redis.DialURL("redis://")
		if err != nil {
			return nil, err
		}
		return c, err
	},
}
// initialize celery client
cli, _ := gocelery.NewCeleryClient(
	gocelery.NewRedisBroker(redisPool),
	&gocelery.RedisCeleryBackend{Pool: redisPool},
	1,
)
// prepare arguments
taskName := "worker.add"
argA := rand.Intn(10)
argB := rand.Intn(10)
// run task
asyncResult, err := cli.Delay(taskName, argA, argB)
if err != nil {
	panic(err)
}
// get results from backend with timeout
res, err := asyncResult.Get(10 * time.Second)
if err != nil {
	panic(err)
}
log.Printf("result: %+v of type %+v", res, reflect.TypeOf(res))
```

## Sample Celery Task Message

Celery Message Protocol Version 1

```javascript
{
    "expires": null,
    "utc": true,
    "args": [5456, 2878],
    "chord": null,
    "callbacks": null,
    "errbacks": null,
    "taskset": null,
    "id": "c8535050-68f1-4e18-9f32-f52f1aab6d9b",
    "retries": 0,
    "task": "worker.add",
    "timelimit": [null, null],
    "eta": null,
    "kwargs": {}
}
```


## Go Celery Worker in Action

![demo](https://raw.githubusercontent.com/gocelery/gocelery/master/demo.gif)

