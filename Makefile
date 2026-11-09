install-swag:
	go install github.com/swaggo/swag/cmd/swag@v1.16.4

install-mockery:
	go install github.com/vektra/mockery/v2@v2.35.4

docs:
	swag init
run:
	go run main.go

ENV ?= local
migrate-up:
	go run main.go -env $(ENV) migrate up
migrate-down:
	go run main.go -env $(ENV) migrate down $(STEPS)
migrate-status:
	go run main.go -env $(ENV) migrate status
seed:
	go run main.go -env $(ENV) seed $(NAME)

mocks:
	mockery --all --keeptree --dir=app/repository --output=app/repository/mocks --case underscore
	mockery --all --keeptree --dir=app/service --output=app/service/mocks --case underscore
	mockery --all --keeptree --dir=app/controller/handler --output=app/controller/handler/mocks --case underscore

test:
	go test -v -coverprofile cover.out ./internal/...
	go tool cover -html cover.out -o cover.html 
