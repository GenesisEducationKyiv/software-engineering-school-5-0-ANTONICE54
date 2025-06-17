# Testing Guide


## Unit Tests

Для запуску unit тестів, використовуйте:

```bash
go test ./internal/...
```

---

## Integration Tests

Для запуску інтеграційних тестів, використовуйте:

```bash
go test ./tests/integration/...
```

---

## End-to-End (E2E) Tests

### Для запуску End-to-End тестів необхідно:

### 1. Запустити сервіс через Docker Compose:

```bash
docker compose up --build -d
```


### 2. Запустити E2E тести:

```bash
go test ./tests/e2e/...
```

---

## General

### Запуск усіх тестів з детальним виводом:

  ```bash
  go test -v ./...
  ```


