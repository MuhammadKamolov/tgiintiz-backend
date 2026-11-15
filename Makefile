# Автоматически подхватываем переменные из .env, если файл существует
ifneq (,$(wildcard ./.env))
	include .env
	export
endif

DB_USER ?= postgres
DB_PASSWORD ?= 123
DB_NAME ?= tgiintiz_db
DB_PORT ?= 5432
DB_URL := postgres://$(DB_USER):$(DB_PASSWORD)@localhost:$(DB_PORT)/$(DB_NAME)?sslmode=disable
MIGRATION_DIR := migrations

.PHONY: run build test sqlc clean \
        migrate-up migrate-down migrate-status \
        infra-up infra-down infra-reset \
        docker-app docker-up docker-down docker-logs docker-ps \
        dev

# === Локальная разработка (Go) ===

run:
	go run ./cmd/api

build:
	go build -ldflags="-s -w" -o bin/api ./cmd/api

test:
	go test -v ./...

sqlc:
	sqlc generate

clean:
	rm -rf bin/

# === Миграции базы данных (Goose) ===

migrate-up:
	goose -dir $(MIGRATION_DIR) postgres "$(DB_URL)" up

migrate-down:
	goose -dir $(MIGRATION_DIR) postgres "$(DB_URL)" down

migrate-status:
	goose -dir $(MIGRATION_DIR) postgres "$(DB_URL)" status

# === Инфраструктура в Docker (только БД и Redis для локальной разработки) ===

infra-up:
	docker compose up -d --wait postgres redis

infra-down:
	docker compose stop postgres redis

infra-reset:
	docker compose down -v
	docker compose up -d --wait postgres redis
	$(MAKE) migrate-up

# === Полный стек в Docker (Production-режим: Postgres + Redis + API) ===

docker-app:
	open -a Docker

docker-up:
	docker compose up -d --build --wait
	$(MAKE) migrate-up

docker-down:
	docker compose down

docker-logs:
	docker compose logs -f

docker-ps:
	docker compose ps

# === Быстрый старт одной командой для разработки ===

dev: infra-up migrate-up run