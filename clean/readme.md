## Тестирование

## Моки

Для генерации используется [mockery](https://vektra.github.io/mockery/latest/).

Запуск генерации:

```sh
mockery
```

Для включения генерации нового интерфейса, дополнить конфиг `.mockery.yaml`

## Проверка кода

## Линтер

Используется [golangci-lint](https://golangci-lint.run/)

Запуск линтера с проверкой синтаксиса:

```sh
go vet ./... && golangci-lint run -v -j $(( $(nproc) - 1))
```

## Тесты

Запустить тесты:

```sh
go test ./...
```

