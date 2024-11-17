install:
	go mod tidy
	go install github.com/rubenv/sql-migrate/...@latest

migrate:
	sql-migrate up -config ./config/sql-migrate.yml -env migration

unmigrate:
	sql-migrate down -config ./config/sql-migrate.yml -env migration

unmigrate-all:
	sql-migrate down -config ./config/sql-migrate.yml -env migration --limit 0

seed:
	sql-migrate up -config ./config/sql-migrate.yml -env seed

dev:
	go run cmd/main.go