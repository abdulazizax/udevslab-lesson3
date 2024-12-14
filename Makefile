run :
	go run cmd/main.go

swag_init:
	swag init -g internal/http/app/app.go --parseDependency -o internal/http/app/docs