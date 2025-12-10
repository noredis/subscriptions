# Subscriptions Service

REST API сервис для агрегации и управления данными об онлайн подписках пользователей.

## 🏗️ Архитектура

Проект реализован с использованием **Onion Architecture** (луковичная архитектура).

### Слои архитектуры
```
┌─────────────────────────────────────┐
│     Presentation Layer (HTTP)       │  ← Fiber handlers, middlewares
├─────────────────────────────────────┤
│       Application Layer             │  ← Application Services, DTOs
├─────────────────────────────────────┤
│         Domain Layer                │  ← Entities, Interfaces, Domain Services, Domain Errors
├─────────────────────────────────────┤
│     Infrastructure Layer            │  ← PostgreSQL, External APIs
└─────────────────────────────────────┘
```

## 🚀 Технологический стек

- **Go 1.21+** - язык программирования
- **Fiber v2** - веб-фреймворк
- **PostgreSQL** - база данных
- **Docker & Docker Compose** - контейнеризация
- **go-playground/validator/v10** - валидация данных
- **golang-migrate** - миграции БД
- **zerolog** - логирование
- **envconfig** - конфигурирование
- **swaggo** - генерация OpenAPI документации

## 📋 Возможности

- Управление подписками пользователей (CRUDL операции)
- Подсчёт суммарной стоимости всех подписок за выбранный период
- Валидация входных данных
- RESTful API с JSON форматом

## 🛠️ Установка и запуск

### Требования

- Docker
- Docker Compose
- Make (опционально)

### Быстрый старт

1. Клонируйте репозиторий:
```bash
git clone https://github.com/noredis/subscriptions.git
cd subscriptions
```

2. Создайте файл `.env`:
```bash
cp .env.example .env
```

3. Настройте переменные окружения в `.env`:
```env
APP_PORT=8080
APP_ENV=dev # dev/prod

LOG_LEVEL=debug # trace/debug/info/warn/error/fatal/panic

DB_USER=postgres
DB_PASSWORD=root
DB_HOST=subs-db
DB_PORT=5432
DB_NAME=subs_db
DB_MAX_CONNS=4
DB_MIN_CONNS=1
DB_MAX_CONN_LIFETIME=1h
DB_MAX_CONN_IDLE_TIME=15m
DB_CONN_ATTEMPTS=5
DB_CONN_DELAY=3s
```

4. Запустите сервис:
```bash
docker-compose up -d
```

5. Проверьте работоспособность:
```bash
curl http://localhost:8080/heartbeat
```

6. Открыть документацию:
```bash
http://localhost:8080/swagger/index.html
```

## 📁 Структура проекта
```
.
├── cmd/
│   └── app/
│       └── main.go                 # Точка входа
├── internal/
│   ├── common/                     # Common Layer
│   │   ├── config/                 # Конфигурация
│   ├── domain/                     # Domain Layer
│   │   ├── entity/                 # Бизнес-сущности
│   │   ├── interfaces/             # Интерфейсы репозиториев
|   |   ├── failure/                # Доменные ошибки
│   │   └── service/                # Доменные сервисы
│   ├── application/                # Application Layer
│   │   ├── appservice/            # Application-сервисы
│   │   └── dto/                    # Data Transfer Objects
│   ├── infrastructure/             # Infrastructure Layer
│   │   └── repository/             # Реализация репозиториев
│   └── presentation/               # Presentation Layer
│       ├── http/
│       │   ├── handlers/           # HTTP обработчики
│       │   └── middlewares/        # Middleware
├── pkg/                            # Shared зависимости
├── migrations/                     # SQL миграции
├── docker-compose.yml
├── Dockerfile
├── Makefile
├── .env
├── go.sum
├── go.mod
└── README.md
```

## 🧪 Тестирование
```bash
# Запустить линтер
make lint
```

## 🔧 Команды Makefile
```bash
make run           # Запустить Docker Compose
make brun          # Пересобрать и запустить Docker Compose
make app-logs      # Показать логи приложения
make lint          # Запустить golangci-lint
make stop          # Остановить контейнеры
make down          # Остановить и удалить контейнеры
make vdown         # Остановить контейнеры и удалить volumes
make db            # Подключиться к PostgreSQL через psql
```

## 🗄️ База данных

### Миграции

Миграции находятся в директории `migrations/` и применяются автоматически при запуске.

### Схема БД

Основные таблицы:
- `subscriptions` - подписки
