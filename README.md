# API Магазин мерча
ТЗ по ссылке https://github.com/avito-tech/tech-internship/tree/main/Tech%20Internships/Backend/Backend-trainee-assignment-winter-2025

Данный сервис реализован на языке Go с использованием библиотеки Gin. Для работы с PostgreSQL м для миграций использовался драйвер pgx. 
Для миграции использовалась библиотека goose, миграции выполняются автоматически при запуске сервиса. Тип операции зависит от определенного запроса. Для данных из файла конфигурации используется Viper.
Реализован магазин мерча, в котором можно покупать товары и передавать монеты, а также получать всю информацию по операциям пользователя. Для токенов авторизации использовалось JWT
В Makefile прописаны возможные варианты запуска API и миграции. Приложение покрыто логами для информировани и дебага.
Для логгирования использовалась Zerolog. Для интеграционного и unit тестирования использовались библиотеки и инструменты testcontainers-go, mockgen, go-sqlmock.
Из дополнительных заданий были реализованы:
- Реализовать интеграционное или E2E-тестирование для остальных сценариев.
- Описать конфигурацию линтера .golangci.yaml в корне проекта для go.
## Запуск приложения:
### Использование docker-compose.
   Для сборки и запуска приложения нужно ввести в консоль команду
   ```
   make build
   ```
   или эту команду
   ```
   docker-compose up --build
   ```
## Пользование сервисом
### 1. Авторизация и регистрация
#### Для отдельной регистрации необходимо выполнить запрос
```
curl --location  --request POST 'http://localhost:8080/api/auth/sign-up' \
--header 'Content-Type: application/json' \
--data '{
    "username": "{username}",
    "password": "{password}"
}'
```
Вместо username вводится желаемый username, в поле password соответственно желаемый пароль.
#### Для авторизации необходимо выполнить запрос
```
curl --location  --request POST 'http://localhost:8080/api/auth/sign-in' \
--header 'Content-Type: application/json' \
--data '{
    "username": "{username}",
    "password": "{password}"
}'
```
Вместо username вводится выбранный нами при регистрации username, в поле password соответственно пароль. Если пользователь не зарегистрирован, то в этом запросе будет сразу осуществлена регистрация и выдача токена.
В ответ на данный запрос нам выдастся токен, который нужно сохранить и использовать во всех следующих запросах. В программе Postman имеется функционал, который позволяет один раз указать токен и выполнять все дальнейшие запросы уже с ним. В командной строке с каждым запросом придется указывать вручную заголовок.
Проверка токена в сервисе выполняется при помощи методов в Middleware.
Во всех запросах вместо Token в заголовке вводится личный токен, полученный при авторизации. 
### 2. Магазин
#### Для покупки мерча необходимо выполнить запрос
```
curl --location --request PUT 'http://localhost:8080/api/buy/{name}' \
--header 'Authorization: Bearer {token}' \
--header 'Content-Type: application/json' \
--data ''
```
Вместо name нужно ввести название желаемого товара для покупки из таблицы, после чего будет выведен id покупки:
| Название     | Цена |
|--------------|------|
| t-shirt      | 80   |
| cup          | 20   |
| book         | 50   |
| pen          | 10   |
| powerbank    | 200  |
| hoody        | 300  |
| umbrella     | 200  |
| socks        | 10   |
| wallet       | 50   |
| pink-hoody   | 500  |

#### Для отправки монет другому пользователю необходимо выполнить запрос
```
curl --location --request POST 'http://localhost:8080/api/sendCoin' \
--header 'Content-Type: application/json' \
--header 'Authorization:  Bearer {token}' \
--data '{
    "destination_username": "{user}",
    "amount": {100}
}'
```
Вместо user в поле нужно ввести никнейм пользователя, которому нужно отправить монеты, а в поле amount количество монет. После успешнего выполнения запроса будет выведен id транзакции.
#### Для получения сгруппированной информации о пользователе необходимо выполнить запрос
```
curl --location 'http://localhost:8080/api/info' \
--header 'Authorization: Bearer {token}' \
--data ''
```
После успешнего выполнения запроса будет выведно количество монет, список купленных им мерчовых товаров и сгруппированная информация о перемещении монеток в кошельке, включая:
- Кто передавал монетки пользователю и в каком количестве
- Кому пользователь передавал монетки и в каком количестве
## Тестирование
Для запуска тестов необходимо ввести команду
```
go test -v ./...
```
Будет выведен процесс результат выполнения всех unit и интеграционных тестов.

Для краткого отображения необходимо ввести команду
```
go test ./...
```
Для работы линтер необходимо ввести
```
golangci-lint run
```
Для работы интеграционных тестов необходим запущенный Docker
## Обработка ошибок
Для различных методов и вызовов функций реализована обработка ошибок, в зависимости от категории ошибки, выдается текст и код ошибки.
