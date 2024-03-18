```docker compose up``` - только БД 

```go test ./... -coverprofile=coverage.out```

```go tool cover -func=coverage.out``` = 54.4% - из-за отсутствия тестов БД

http://localhost:8082/swagger/index.html#/

Required env CONFIG_PATH={YOUR_PATH}/film_library/config/local.yaml
