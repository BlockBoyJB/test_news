# test_news

### Prerequisites

- Docker, Docker Compose
- or Golang 1.25 + postgresql

### Getting started

* Добавить репозиторий к себе
* Создать .env файл в директории с проектом и заполнить информацией из .env.example

### Usage

Запустить сервис можно с помощью `make compose-up` (или `docker-compose up -d --build`)
или `make run` (при наличии go1.25 и локально развернутого postgresql)  
Тесты доступны по команде `make tests`

### Примеры запросов

#### Получение токена

Для работы с api нужна авторизация через Bearer токен и Authorization заголовок  
`request`

```shell
curl -X 'GET' \
  'http://localhost:8000/authorize'
```

`response`

```shell
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NjE5MDc4MDYsImlhdCI6MTc2MTMwMzAwNn0.wgmhVaaWqXRbAM8GTHdUjwvjLgsYPPok4UTrOyq78Zg"
}
```

#### Создание новости

Требуется токен
`request`

```shell
curl -X 'POST' \
  'http://localhost:8000/api/v1/news/create' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NjE5MDc4MDYsImlhdCI6MTc2MTMwMzAwNn0.wgmhVaaWqXRbAM8GTHdUjwvjLgsYPPok4UTrOyq78Zg' \
  -d '{"Title": "hello world", "Content": "my content", "Categories": [1, 2, 3]}'
```

`response`  
В ответе - id созданной новости

```json
{
  "Id": 1
}
```

#### Обновление новости

Обновление новости по id (передается в path). Все параметры - опциональные (если не поле указано, то в БД не обновится).
Если указано поле Categories, то происходит замена всех существующих на новые
`request`

```shell
curl -X 'POST' \
  'http://localhost:8000/api/v1/news/edit/1' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NjE5MDc4MDYsImlhdCI6MTc2MTMwMzAwNn0.wgmhVaaWqXRbAM8GTHdUjwvjLgsYPPok4UTrOyq78Zg' \
  -d '{"Title": "New TITLE", "Categories": [2]}'
```

`response`
`OK`

#### Получение списка новостей с пагинацией

limit, offset пагинация (для limit допустимый диапазон - [0, 20])

`request`

```shell
curl -X 'GET' \
  'http://localhost:8000/api/v1/news/list?limit=10&offset=0' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NjE5MDc4MDYsImlhdCI6MTc2MTMwMzAwNn0.wgmhVaaWqXRbAM8GTHdUjwvjLgsYPPok4UTrOyq78Zg'
```

`response`

```json
{
  "Success": true,
  "News": [
    {
      "Id": 1,
      "Title": "New TITLE",
      "Content": "my content",
      "Categories": [
        2
      ]
    }
  ]
}

```