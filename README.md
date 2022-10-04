# Сервис для подачи объявлений
![Project language][badge_language]
[![Test & Lint Status][badge_build]][link_build]
[![codecov](https://codecov.io/gh/nizhikebinesi/golang-test-task/graph/badge.svg?token=JJVKAZ8PWX)](https://codecov.io/gh/nizhikebinesi/golang-test-task)
[![Twitter Follow](https://img.shields.io/twitter/follow/nizhikebinesi)](https://twitter.com/nizhikebinesi)


[badge_build]:https://img.shields.io/github/workflow/status/nizhikebinesi/golang-test-task/Check%20on%20PRs%20and%20push
[badge_language]:https://img.shields.io/badge/language-go_1.18-blue.svg?longCache=true
[link_build]:https://github.com/nizhikebinesi/golang-test-task/actions


## Как запустить
0. Установить Docker и [Docker-Compose](https://www.digitalocean.com/community/tutorials/how-to-install-and-use-docker-compose-on-ubuntu-20-04-ru)
1. `docker-compose build && docker-compose up -d`
2. Таблицы создаются через `gorm` внутри `app`, БД создается пустой.

## API
Address: `http://localhost:8888`
Prefix: `/api/v0.1`

| Endpoint    | Method | Description                                |
|-------------| ------------- |--------------------------------------------|
| `/ads`      | `POST` | Создание объявления |                       
| `/ads`      | `GET` | Получение списка объявлений(для пагинции)  |
| `/ads/{id}` | `GET` | Получение объявления |                      

## Примеры использования
### 1. Создать объявление(`/create_ad`)

Пример запроса:
```shell
curl POST -v -d "{
  \"title\": \"Гараж\",
  \"description\": \"Продам гараж\nТелефон:...\",
  \"price\": 45999.99,
  \"image_urls\": [\"photo.example.ru/img/1.png\", \"photo.example.ru/img/2.jpg\", \"photo.example.ru/img/3.jpeg\"],
}" http://localhost:8888/v0.1/create_ad
```
Пример успешного ответа:
```json
{
  "status": "success",
  "id": 1
}
```
Пример неудачного ответа:

```json
{
  "status": "error",
  "id": null
}
```

### 2. Получить объявление по ID(`/get_ad`)
Пример запроса:
```shell
curl -X GET "http://localhost:8888/v0.1/get_ad/1?fields=description&fields=image_urls"
```
Примеры успешных ответов:
```json
{
  "status": "success",
  "result": {
    "id": 1,
    "title": "Гараж",
    "description": "Продам гараж!",
    "price": 99999.99,
    "main_image_url": "garage.ru/img/1.png",
    "image_urls": ["supergarage.ru/img/2.jpeg"]
  }
}
```
Или:
```json
{
  "status": "success",
  "result": {
    "id": 1,
    "title": "Гараж",
    "description": "Продам гараж!",
    "price": 99999.99,
    "main_image_url": null,
    "image_urls": []
  }
}
```

Пример неудачных ответов:
```json
{
  "status": "error",
  "result": null
}
```

### 3.1. Получить список объявлений для пагинации(вторые 10 объявлений с сортировкой по дате по возрастанию, в случае 13 объявлений в БД)
Пример запроса:
```shell
curl -X GET "http://localhost:8888/v0.1/list_ads?by=created_at&asc=true&offset=10"
```
Пример ответа:
```json
{
  "status": "success",
  "result": [
    {
      "id": 11,
      "title": "Гараж",
      "description": "Продам гараж!",
      "price": 99999.99,
      "main_image_url": null,
      "image_urls": []
    },
    {
      "id": 12,
      "title": "Тетрадь в клетку",
      "description": "",
      "price": 79.89,
      "main_image_url": null,
      "image_urls": []
    },
    {
      "id": 13,
      "title": "Книга №42",
      "description": "",
      "price": 979.29,
      "main_image_url": null,
      "image_urls": []
    }
  ]
}
```

### 3.2. Получить список объявлений для пагинации(первые 10 объявлений с сортировкой по цене по убыванию)
Пример запроса:
```shell
curl -X GET "http://localhost:8888/v0.1/list_ads?by=price&asc=false"
```
Пример ответа:
```json
{
  "status": "success",
  "result": [
    {
      "id": 1,
      "title": "Гараж",
      "description": "Продам гараж!",
      "price": 99999.99,
      "main_image_url": "garage.ru/img/1.png",
      "image_urls": []
    },
    {
      "id": 2,
      "title": "Тетрадь в клетку",
      "description": "",
      "price": 79.89,
      "main_image_url": null,
      "image_urls": []
    },
    {
      "id": 3,
      "title": "Книга №42",
      "description": "",
      "price": 979.29,
      "main_image_url": null,
      "image_urls": []
    },
    {
      "id": 4,
      "title": "Уловка 22",
      "description": "Роман американского писателя Джозефа Хеллера, опубликованный в 1961 году. ",
      "price": 629.99,
      "main_image_url": null,
      "image_urls": []
    },
    {
      "id": 5,
      "title": "Обои",
      "description": "Рулон обоев",
      "price": 2629.99,
      "main_image_url": null,
      "image_urls": []
    },
    {
      "id": 6,
      "title": "Воздух",
      "description": "Свежий воздух",
      "price": 9999999.99,
      "main_image_url": null,
      "image_urls": []
    },
    {
      "id": 7,
      "title": "Дженга",
      "description": "Настольная игра...",
      "price": 999,
      "main_image_url": null,
      "image_urls": []
    },
    {
      "id": 8,
      "title": "Сахар",
      "description": "На развес, цена за 1 кг",
      "price": 57.87,
      "main_image_url": null,
      "image_urls": []
    },
    {
      "id": 9,
      "title": "Соль",
      "description": "Цена за килограмм",
      "price": 15,
      "main_image_url": null,
      "image_urls": []
    },
    {
      "id": 10,
      "title": "Ноутбук",
      "description": "",
      "price": 16999.99,
      "main_image_url": "top.laptops/ultra/1.jpeg",
      "image_urls": []
    }
  ]
}
```

## TODOs
1. [x] Добавить `Sentry`
2. [ ] Добавить `Prometheus`, `Grafana`
3. [ ] Добавить `Master-Slave репликацию` для `Postgres`
4. [ ] Добавить `HA Proxy`/`Consul`/`pgpool`/`pgbouncer`
5. [ ] Добавить `nginx`(с `Consul`) и запустить копии сервиса
6. [ ] Сгенерировать и хостить `Swagger`-документацию
7. [ ] Добавить `DELETE`(удаления записей) и `PUT`(изменения записей) методы в сервис
8. [ ] Добавить тестирование через `dockertest`

## Оригинальный текст
### **Задача**

Разработать сервис для подачи объявлений с сохранением в базе данных. 
Сервис должен предоставлять API, работающее поверх HTTP в формате JSON.

### **Требования**

- Язык программирования — Go;
- Готовую версию выложить на Github;
- Простая инструкция для запуска(в идеале 
— с возможностью запустить через `docker-compose up`, но это необязательно);
- 3 метода:
    - получение списка объявлений,
    - получение одного объявления,
    - создание объявления;
- Валидация полей:
    - не больше 3 ссылок на фото,
    - описание не больше 1000 символов,
    - название не больше 200 символов;

Если есть сомнения по деталям — решение принять самостоятельно, 
но в своём `README.md` рекомендуем выписать вопросы и принятые решения по ним.

### Ограничения по времени

2-4 часа на выполнение. Если что-то не укладывается в указанное время, 
то реализовать задачу по степени важности функционала. 
Мы не требуем выполнить абсолютно всё. Здесь важны умение приоритизировать и 
чистота кода.

### **Детали**

**Метод получения списка объявлений**

- Пагинация: на одной странице должно присутствовать 10 объявлений;
- Cортировки: по цене(возрастание/убывание) и 
по дате создания(возрастание/убывание);
- Поля в ответе: название объявления, 
ссылка на главное фото (первое в списке), цена.

**Метод получения конкретного объявления**

- Обязательные поля в ответе: название объявления, цена, ссылка на главное фото;
- Опциональные поля (можно запросить, передав параметр fields): 
описание, ссылки на все фото.

**Метод создания объявления:**

- Принимает все вышеперечисленные поля: название, описание, 
несколько ссылок на фотографии
(сами фото загружать никуда не требуется), цена;
- Возвращает ID созданного объявления и код результата (ошибка или успех).
