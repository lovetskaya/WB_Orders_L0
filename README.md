## Техническое задание 
Необходимо разработать демонстрационный сервис с простейшим интерфейсом, отображающий данные о заказе. Модель данных можно найти в конца задания.
Сервис устанавливает соединение с базой данных PostgreSQL, подписывается на Kafka для получения информации о заказах, кэширует данные в оперативной памяти и предоставляет HTTP API для доступа к информации о заказах по их идентификаторам.

Основные аспекты:
1. Развернуть локально PostgreSQL
2. Создать свою БД
3. Настроить своего пользователя
4. Создать таблицы для хранения полученных данных
5. Разработать сервис
6. Реализовать подключение к брокерам и подписку на топик orders в Kafka
7. Полученные данные записывать в БД
8. Реализовать кэширование полученных данных в сервисе (сохранять in memory)
9. В случае прекращения работы сервиса необходимо восстанавливать кэш из БД
10. Запустить http-сервер и выдавать данные по id из кэша
11. Разработать простейший интерфейс отображения полученных данных по id заказа

Дополнительные аспекты:
1. Сделать отдельный скрипт для публикации данных в топик, чтобы проверить работает ли подписка онлайн
2. Реализовать сохранение данные в случае ошибок или проблем с сервисом
3. Удобно разворачивать сервисы для тестирования можно через docker-compose
4. Покройте сервис автотестами
5. Устройте вашему сервису стресс-тест: выясните на что он способен (утилиты WRK и Vegeta)
6. Логи в JSON-формате делают их более удобными для машинной обработки

## Требования
Go 1.22.3
PostgreSQL
Kafka & Zookeper
WRK/Vegeta
Docker, docker-compose

## Развертывание
1. Клонирование репозитория

```git clone https://github.com/lovetskaya/WB_Orders_L0.git```

2. Переход в рабочую директорию

```cd wb_service_order```

3. Установка зависимостей

```go mod tidy```

4. Создайте файл конфигурации в корневой директории проекта с именем config.docker.yaml и опишите его [вот так](docker-compose.yml)

## Настройка базы данных

Создайте БД wb_orders_db и далее выполните начальную миграцию с помощью скрипта SQL:

```
CREATE DATABASE wb_orders_db;
CREATE TABLE orders (
       order_uid VARCHAR(50) PRIMARY KEY,
       track_number VARCHAR(50),
       entry VARCHAR(10),
       delivery_name VARCHAR(100),
       delivery_phone VARCHAR(20),
       delivery_zip VARCHAR(10),
       delivery_city VARCHAR(100),
       delivery_address VARCHAR(255),
       delivery_region VARCHAR(100),
       delivery_email VARCHAR(100),
       payment_transaction VARCHAR(50),
       payment_request_id VARCHAR(50),
       payment_currency VARCHAR(3),
       payment_provider VARCHAR(50),
       payment_amount DECIMAL,
       payment_payment_dt TIMESTAMP,
       payment_bank VARCHAR(50),
       payment_delivery_cost DECIMAL,
       payment_goods_total DECIMAL,
       payment_custom_fee DECIMAL,
       items JSONB,
       locale VARCHAR(10),
       internal_signature TEXT,
       customer_id VARCHAR(50),
       delivery_service VARCHAR(50),
       shardkey VARCHAR(10),
       sm_id INT,
       date_created TIMESTAMP,
       oof_shard VARCHAR(10)
   );
```
## Данные о заказе доступны по адресу:

[[http://localhost:8080/order?id=b563feb7b2b84b6test](http://localhost:8080/order?id=b563feb7b2b84b6test)]

В параметр id GET-запроса необходимо подставить ID требующегося заказа. При переходе по ссылке выше отобразится информация о заказе с id = b563feb7b2b84b6test.

## Запуск сервиса в Docker
Соберите Docker-образы с использованием переменной окружения APP_ENV=docker и запустите конейнеры в фоновом режиме:

```make docker```

Запустите сервис:

```make run```

## Запуск тестов
1. Нагрузочные тесты:
Тест проводит нагрузочное тестирование, используя 12 потоков (-t12) и поддерживая 400 соединений (-c400) в течение 30 секунд (-d30s). Он отправляет запросы к указанному URL:

```make wrk-docker```

Тест выполняет нагрузочное тестирование, отправляя 10 запросов в секунду в течение 30 секунд, используя цели, указанные в файле targets.txt. Результаты теста сохраняются в файл results.bin:

```make vegeta-docker```

2. Автотест (проверка обработчика GetOrder для получения информации о заказе через mock-имитацию поведения БД):

```make test-docker```


