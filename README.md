# raketa-backend

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