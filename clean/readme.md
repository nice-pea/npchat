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