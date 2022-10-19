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

## Go Celery Worker in Action

![demo](https://github.com/sQUARys/TestTaskMailGaner/blob/master/project.gif)

