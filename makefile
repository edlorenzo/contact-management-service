#!make
include .env
export $(shell sed 's/=.*//' .env)

run:
	go run cmd/main.go

migrate-up:
	migrate -database "postgresql://$(DB_USERNAME):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_DATABASE)?sslmode=disable" -path "cmd/migrations/" up

migrate-down:
	migrate -database "postgresql://$(DB_USERNAME):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_DATABASE)?sslmode=disable" -path "cmd/migrations/" down

create-migration:
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir cmd/migrations -seq $$name