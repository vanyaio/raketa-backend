# raketa-backend

# Перед запуском проекта локально или через docker-compose в корне создайте файл .env:
```
    ADMIN_RAKETA=qwerty
    GRPC_PORT=:50052
    REST_PORT=:9090
    POSTGRES_USER=postgres
    POSTGRES_PASSWORD=postgres
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

# SWAGGER доступен после запуска проекта "make run"
```
    http://localhost:9090/swagger
```

# Для теститорования методов storage:
```
    make postgres-local-up - если базы нет
    make postgres-local-run - если база есть
    make migrate-up -если миграции не сделаны
    make test - для самого тестирования
```
