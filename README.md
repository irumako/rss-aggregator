# RSS Aggregator

RSS-агрегатор на Go с REST API для управления RSS-лентами.

## Требования

- **Go** 1.25.0 или выше
- **oapi-codegen** v2.5.1+ (для генерации кода из OpenAPI спецификации)
- **goose** (для миграций базы данных)

## Установка зависимостей

```bash
go mod download
```

## Генерация кода

Проект использует генерацию кода из OpenAPI спецификации (`api/swagger.yml`) в Go-код.

### Установка oapi-codegen

```bash
go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest
```

### Генерация кода

Для генерации кода из `api/swagger.yml` используйте следующую команду:

```bash
oapi-codegen -config config/gen.yml api/swagger.yml
```

Эта команда:
- Читает конфигурацию из `config/gen.yml`
- Генерирует Fiber-сервер и модели из `api/swagger.yml`
- Сохраняет результат в `gen/gen.go`

**Важно:** Не редактируйте файл `gen/gen.go` вручную — он автоматически генерируется.

## Миграции базы данных

Проект использует [goose](https://github.com/pressly/goose) для управления миграциями.

### Установка goose

```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```

### Применение миграций

По умолчанию используется SQLite3 с базой данных `./rss.db`:

```bash
make migrate-up
```

Или с явным указанием строки подключения:

```bash
make migrate-up DB_STRING=sqlite3://./rss.db
```

### Другие команды миграций

```bash
# Показать статус миграций
make migrate-status

# Откатить последнюю миграцию
make migrate-down

# Откатить миграции до указанной версии
make migrate-down-to VERSION=20240101000000

# Применить миграции до указанной версии
make migrate-up-to VERSION=20240101000000

# Создать новую миграцию
make migrate-create NAME=add_new_table

# Показать справку по всем командам
make help
```

## Тестирование

### Запуск всех тестов

```bash
go test ./...
```

### Запуск интеграционных тестов

```bash
go test ./internal/service/... -v
```

### Запуск тестов с покрытием

```bash
go test ./... -cover
```

## Сборка и запуск

### Сборка сервера

```bash
go build -o server.exe ./cmd/server
```

### Запуск сервера

```bash
./server.exe
```

Или с переменными окружения:

```bash
DB_PATH=./rss.db PORT=3000 ./server.exe
```

По умолчанию:
- База данных: `./rss.db`
- Порт: `3000`

## Проверка работоспособности

### 1. Проверка генерации кода

После изменения `api/swagger.yml` необходимо перегенерировать код:

```bash
# Генерация кода
oapi-codegen -config config/gen.yml api/swagger.yml

# Проверка, что код компилируется
go build ./...
```


### 2. Проверка тестов

```bash
# Запустить все тесты
go test ./...

# Запустить интеграционные тесты
go test ./internal/service/... -v
```

### 3. Проверка сборки

```bash
# Собрать проект
go build -o server.exe ./cmd/server

# Запустить сервер
./server.exe
```

