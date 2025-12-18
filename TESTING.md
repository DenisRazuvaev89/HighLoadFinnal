# Инструкция по тестированию

## 1. Локальный запуск

```bash
go build -o microservice .
./microservice
```

Сервис запустится на порту 8081.

## 2. Базовое тестирование API

### Создание пользователя
```bash
curl -X POST http://localhost:8081/api/users \
  -H "Content-Type: application/json" \
  -d '{"name":"Test User","email":"test@example.com"}'
```

Ожидаемый результат: `{"id":1,"name":"Test User","email":"test@example.com"}`

### Получение всех пользователей
```bash
curl http://localhost:8081/api/users
```

### Получение пользователя по ID
```bash
curl http://localhost:8081/api/users/1
```

### Обновление пользователя
```bash
curl -X PUT http://localhost:8081/api/users/1 \
  -H "Content-Type: application/json" \
  -d '{"name":"Updated User","email":"updated@example.com"}'
```

### Удаление пользователя
```bash
curl -X DELETE http://localhost:8081/api/users/1
```

### Проверка метрик
```bash
curl http://localhost:8081/metrics
```

### Health check
```bash
curl http://localhost:8081/health
```

## 3. Тестирование валидации

### Невалидный email
```bash
curl -X POST http://localhost:8081/api/users \
  -H "Content-Type: application/json" \
  -d '{"name":"Test","email":"invalid-email"}'
```

Ожидаемый результат: ошибка 400 с сообщением "invalid email format"

### Пустое имя
```bash
curl -X POST http://localhost:8081/api/users \
  -H "Content-Type: application/json" \
  -d '{"name":"","email":"test@example.com"}'
```

Ожидаемый результат: ошибка 400 с сообщением "name is required"

## 4. Нагрузочное тестирование с wrk

Установите wrk:
```bash
brew install wrk
```

Запустите тест:
```bash
wrk -t12 -c500 -d60s http://localhost:8081/api/users
```

Ожидаемые результаты:
- Requests/sec: > 1000
- Latency avg: < 10ms
- Errors: 0

## 5. Docker тестирование

### Сборка образа
```bash
docker build -t go-microservice .
```

### Запуск контейнера
```bash
docker run -p 8081:8081 go-microservice
```

### Запуск через docker-compose
```bash
docker-compose up --build
```

Сервисы:
- Микросервис: http://localhost:8081
- Prometheus: http://localhost:9091

### Остановка
```bash
docker-compose down
```

## 6. Проверка логирования

Запустите несколько операций и проверьте логи:
```bash
curl -X POST http://localhost:8081/api/users -H "Content-Type: application/json" -d '{"name":"User1","email":"user1@example.com"}'
curl -X POST http://localhost:8081/api/users -H "Content-Type: application/json" -d '{"name":"User2","email":"user2@example.com"}'
curl http://localhost:8081/api/users/1
curl -X DELETE http://localhost:8081/api/users/1
```

В консоли должны появиться AUDIT логи:
```
[AUDIT] 2025-12-18T... - Action: CREATE - UserID: 1
[AUDIT] 2025-12-18T... - Action: CREATE - UserID: 2
[AUDIT] 2025-12-18T... - Action: GET - UserID: 1
[AUDIT] 2025-12-18T... - Action: DELETE - UserID: 1
```

## 7. Проверка метрик Prometheus

Запустите несколько запросов и проверьте метрики:
```bash
curl http://localhost:8081/metrics | grep http_requests_total
curl http://localhost:8081/metrics | grep http_request_duration_seconds
```

## 8. Тестирование Rate Limiting

Для проверки rate limiter (1000 req/s с burst 5000) используйте wrk:
```bash
wrk -t20 -c1000 -d10s http://localhost:8081/api/users
```

При превышении лимита должны появиться ошибки 429 (Too Many Requests).

## 9. Проверка конкурентности

Создайте несколько пользователей одновременно:
```bash
for i in {1..10}; do
  curl -X POST http://localhost:8081/api/users \
    -H "Content-Type: application/json" \
    -d "{\"name\":\"User$i\",\"email\":\"user$i@example.com\"}" &
done
wait
```

Проверьте, что все пользователи созданы:
```bash
curl http://localhost:8081/api/users
```

## Чек-лист для сдачи

- [ ] Сервис компилируется без ошибок
- [ ] Все CRUD операции работают корректно
- [ ] Валидация email работает
- [ ] Асинхронное логирование работает (проверить логи)
- [ ] Rate limiting настроен (1000 req/s, burst 5000)
- [ ] Метрики Prometheus доступны на /metrics
- [ ] Health check работает на /health
- [ ] Docker образ собирается успешно
- [ ] docker-compose up запускает оба сервиса
- [ ] wrk показывает > 1000 RPS с latency < 10ms
- [ ] 0 ошибок при нагрузочном тесте
