.PHONY: help migrate-up migrate-down migrate-status migrate-create migrate-reset migrate-reset-confirm migrate-version migrate-down-to migrate-up-to

# Database connection string
# Example: postgres://user:password@localhost:5432/dbname?sslmode=disable
# Example: mysql://user:password@tcp(localhost:3306)/dbname
# Example: sqlite3:///path/to/database.db
DB_STRING ?= ./rss.db

# Database type (auto-detected from DB_STRING if not set)
# Can be: postgres, mysql, sqlite3
DB_TYPE ?= $(if $(findstring postgres://,$(DB_STRING)),postgres,$(if $(findstring mysql://,$(DB_STRING)),mysql,$(if $(findstring sqlite3://,$(DB_STRING)),sqlite3,postgres)))

# Migration directory
MIGRATIONS_DIR ?= ./migrations

# Goose binary path (will use go run if not installed)
GOOSE_CMD ?= go run github.com/pressly/goose/v3/cmd/goose@latest

help: ## Показать справку по доступным командам
	@echo Доступные команды:
	@echo 
	@echo   migrate-up              - Применить все миграции
	@echo   migrate-down            - Откатить последнюю миграцию
	@echo   migrate-down-to         - Откатить миграции до указанной версии (VERSION=20240101000000)
	@echo   migrate-up-to           - Применить миграции до указанной версии (VERSION=20240101000000)
	@echo   migrate-status          - Показать статус миграций
	@echo   migrate-version         - Показать текущую версию миграции
	@echo   migrate-create          - Создать новую миграцию (NAME=create_users_table)
	@echo   migrate-reset           - Откатить все миграции и применить заново (ОПАСНО!)
	@echo   migrate-reset-confirm   - Подтвержденный сброс миграций (ОПАСНО!)
	@echo   migrate-up-mysql        - Применить все миграции (MySQL)
	@echo   migrate-up-sqlite3      - Применить все миграции (SQLite3)
	@echo   migrate-status-mysql    - Показать статус миграций (MySQL)
	@echo   migrate-status-sqlite3  - Показать статус миграций (SQLite3)
	@echo 
	@echo Переменные:
	@echo   DB_STRING               - Строка подключения к БД (по умолчанию: sqlite3://rss.db)
	@echo   MIGRATIONS_DIR          - Директория с миграциями (по умолчанию: migrations)
	@echo 
	@echo Примеры:
	@echo   make migrate-create NAME=init
	@echo   make migrate-up DB_STRING="postgres://user:pass@localhost:5432/dbname?sslmode=disable"

migrate-up: ## Применить все миграции
	$(GOOSE_CMD) -dir $(MIGRATIONS_DIR) $(DB_TYPE) "$(DB_STRING)" up

migrate-down: ## Откатить последнюю миграцию
	$(GOOSE_CMD) -dir $(MIGRATIONS_DIR) $(DB_TYPE) "$(DB_STRING)" down

migrate-down-to: ## Откатить миграции до указанной версии (использование: make migrate-down-to VERSION=20240101000000)
	$(if $(VERSION),,$(error VERSION is required. Example: make migrate-down-to VERSION=20240101000000))
	$(GOOSE_CMD) -dir $(MIGRATIONS_DIR) $(DB_TYPE) "$(DB_STRING)" down-to $(VERSION)

migrate-up-to: ## Применить миграции до указанной версии (использование: make migrate-up-to VERSION=20240101000000)
	$(if $(VERSION),,$(error VERSION is required. Example: make migrate-up-to VERSION=20240101000000))
	$(GOOSE_CMD) -dir $(MIGRATIONS_DIR) $(DB_TYPE) "$(DB_STRING)" up-to $(VERSION)

migrate-status: ## Показать статус миграций
	$(GOOSE_CMD) -dir $(MIGRATIONS_DIR) $(DB_TYPE) "$(DB_STRING)" status

migrate-version: ## Показать текущую версию миграции
	$(GOOSE_CMD) -dir $(MIGRATIONS_DIR) $(DB_TYPE) "$(DB_STRING)" version

migrate-create: ## Создать новую миграцию (использование: make migrate-create NAME=create_users_table)
	$(if $(NAME),,$(error NAME is required. Example: make migrate-create NAME=create_users_table))
	$(GOOSE_CMD) -dir $(MIGRATIONS_DIR) create $(NAME) sql

migrate-reset: ## Откатить все миграции и применить заново (ОПАСНО!)
	@echo "Внимание: это откатит все миграции!"
	@echo "Для подтверждения выполните: make migrate-reset-confirm"
	@exit 1

migrate-reset-confirm: ## Подтвержденный сброс миграций (ОПАСНО!)
	$(GOOSE_CMD) -dir $(MIGRATIONS_DIR) $(DB_TYPE) "$(DB_STRING)" reset
	$(GOOSE_CMD) -dir $(MIGRATIONS_DIR) $(DB_TYPE) "$(DB_STRING)" up

# Альтернативные команды для других СУБД

migrate-up-mysql: ## Применить все миграции (MySQL)
	$(GOOSE_CMD) -dir $(MIGRATIONS_DIR) mysql "$(DB_STRING)" up

migrate-up-sqlite3: ## Применить все миграции (SQLite3)
	$(GOOSE_CMD) -dir $(MIGRATIONS_DIR) sqlite3 "$(DB_STRING)" up

migrate-status-mysql: ## Показать статус миграций (MySQL)
	$(GOOSE_CMD) -dir $(MIGRATIONS_DIR) mysql "$(DB_STRING)" status

migrate-status-sqlite3: ## Показать статус миграций (SQLite3)
	$(GOOSE_CMD) -dir $(MIGRATIONS_DIR) sqlite3 "$(DB_STRING)" status

