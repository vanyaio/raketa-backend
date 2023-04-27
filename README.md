# raketa-backend

# Перед запуском проекта локально или через docker-compose в корне создайте файл .env:
```
    ADMIN_RAKETA=1234
    GRPC_PORT=:50052
    REST_PORT=:9090
    POSTGRES_USER=postgres
    POSTGRES_PASSWORD=postgres
```

# Для теститорования методов storage запустите базу данных, создайте миграции и создайте в папке с тестами файл .env:
```
    POSTGRES_HOST=localhost
    POSTGRES_USER=postgres
    POSTGRES_PASSWORD=postgres
    POSTGRES_DB=raketalocaldb
```
# Далее введите команду:
```
    make tests
```

# Запуск через docker-compose
```
    docker-compose  up --build raketa
```
# Выключение docker-compose
```
    docker-compose down
```
# Для отката миграций в docker-compose
```
    make docker-migrate-down - откат всех миграций
    make docker-migrate-down version=1 - откат до 1-й миграции
```


# Запуск серверов локально
```
    make postgres-up
    make migrate-up
    make run
```

# Удаление базы
```
    make postgres-local-del
```