.PHONY: db
db:
	docker compose up -d db

.PHONY: db-connect
db-connect:
	docker compose exec db psql -U omigoto

.PHONY: db-export
db-export:
	docker compose exec db psqldef -U omigoto omigoto --export > ./db/current.sql

.PHONY: db-migrate
db-migrate:
	docker compose exec -T db psqldef -U omigoto omigoto --enable-drop < ./db/schema.sql

.PHONY: db-migrate-dry-run
db-migrate-dry-run:
	docker compose exec -T db psqldef -U omigoto omigoto --enable-drop --dry-run < ./db/schema.sql

.PHONY: sqlc
sqlc:
	docker run --rm -v $(PWD):/src -w /src sqlc/sqlc generate -f ./backend/db/sqlc.yaml
