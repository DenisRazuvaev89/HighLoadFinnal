# Go Microservice - High Load Service

Микросервис на Go для обработки высоких нагрузок с CRUD операциями для управления пользователями.

## Функциональность

- RESTful API для управления пользователями
- Конкурентная обработка с использованием goroutines
- Rate limiting (1000 req/s с burst 5000)
- Prometheus метрики (RPS, latency)
- Асинхронное логирование операций
- Валидация данных (включая email)
- Docker контейнеризация

## Технологии

- Go 1.22+
- gorilla/mux - HTTP роутинг
- golang.org/x/time/rate - rate limiting
- prometheus/client_golang - метрики
- Docker & Docker Compose

## API Endpoints

| Метод | Endpoint | Описание |
|-------|----------|----------|
| GET | /api/users | Получить всех пользователей |
| GET | /api/users/{id} | Получить пользователя по ID |
| POST | /api/users | Создать нового пользователя |
| PUT | /api/users/{id} | Обновить пользователя |
| DELETE | /api/users/{id} | Удалить пользователя |
| GET | /metrics | Prometheus метрики |
| GET | /health | Health check |

## Быстрый старт

### Локальный запуск

```bash
go mod download
go build -o microservice .
./microservice
```

Сервис будет доступен на http://localhost:8081

### Docker запуск

```bash
docker-compose up --build
```

Сервисы:
- Microservice: http://localhost:8081
- Prometheus: http://localhost:9091

## Примеры использования

### Создать пользователя

```bash
curl -X POST http://localhost:8081/api/users \
  -H "Content-Type: application/json" \
  -d '{"name":"John Doe","email":"john@example.com"}'
```

### Получить всех пользователей

```bash
curl http://localhost:8081/api/users
```

### Получить пользователя по ID

```bash
curl http://localhost:8081/api/users/1
```

### Обновить пользователя

```bash
curl -X PUT http://localhost:8081/api/users/1 \
  -H "Content-Type: application/json" \
  -d '{"name":"Jane Doe","email":"jane@example.com"}'
```

### Удалить пользователя

```bash
curl -X DELETE http://localhost:8081/api/users/1
```

### Получить метрики

```bash
curl http://localhost:8081/metrics
```

## Нагрузочное тестирование

Установите wrk:
```bash
# macOS
brew install wrk

# Linux
sudo apt-get install wrk
```

Запустите тест:
```bash
wrk -t12 -c500 -d60s http://localhost:8081/api/users
```

Ожидаемые результаты:
- RPS: > 1000 запросов в секунду
- Latency: < 10 мс
- Ошибок: 0%

## Архитектура

```
go-microservice/
├── main.go                 # Точка входа
├── handlers/              # HTTP обработчики
│   └── user_handler.go
├── services/              # Бизнес-логика
│   └── user_service.go
├── models/                # Модели данных
│   └── user.go
├── utils/                 # Утилиты
│   ├── logger.go
│   └── rate_limiter.go
├── metrics/               # Prometheus метрики
│   └── prometheus.go
├── Dockerfile
├── docker-compose.yml
└── prometheus.yml
```

## Ключевые особенности реализации

### Rate Limiting
- Настроен лимит 1000 req/s
- Burst 5000 для обработки всплесков нагрузки
- Использует токен-бакет алгоритм

### Конкурентность
- Асинхронное логирование через goroutines
- Thread-safe операции с использованием sync.RWMutex
- Неблокирующая обработка audit log

### Метрики
- http_requests_total - счетчик запросов
- http_request_duration_seconds - гистограмма latency
- Экспорт в формате Prometheus

### Валидация
- Проверка обязательных полей
- Валидация формата email через regex
- Обработка ошибок декодирования JSON

## Мониторинг

Prometheus доступен на http://localhost:9091

Примеры запросов:
- Текущий RPS: `rate(http_requests_total[1m])`
- Средняя latency: `rate(http_request_duration_seconds_sum[1m]) / rate(http_request_duration_seconds_count[1m])`

## Особенности разработки

### Изменения относительно базового кода

1. **Увеличен burst до 5000** - для прохождения теста с 500 одновременными соединениями
2. **Добавлена валидация email** - regex проверка формата
3. **Исправлен ID=0 в логах** - логирование после присвоения ID
4. **Добавлена обработка ошибок** - валидация и error handling во всех handlers
5. **Health check endpoint** - для Docker healthcheck
6. **Prometheus конфигурация** - для автоматического сбора метрик

## Требования к окружению

- Go 1.22 или выше
- Docker и Docker Compose (опционально)
- wrk для нагрузочного тестирования

## Лицензия

Образовательный проект для изучения high-load архитектур.