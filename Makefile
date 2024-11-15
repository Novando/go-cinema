migrate:
	sql-migrate up -config ./config/sql-migrate.yml -env migration

unmigrate:
	sql-migrate down -config ./config/sql-migrate.yml -env migration

unmigrate-all:
	sql-migrate down -config ./config/sql-migrate.yml -env migration --limit 0

coba:
	pwd