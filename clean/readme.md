### Моки

Для генерации используется [mockery](https://vektra.github.io/mockery/latest/).

Запуск генерации:

```shellscript
mockery
```

Для включения генерации нового интерфейса, дополнить конфиг `.mockery.yaml`

### Линтер

Используется [golangci-lint](https://golangci-lint.run/)

Запуск линтера с проверкой синтаксиса:

```shellscript
go vet ./... && golangci-lint run -v -j $(( $(nproc) - 1))
```

### Тесты

Запустить тесты:

```shellscript
go test ./...
```