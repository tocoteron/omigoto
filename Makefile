.PHONY: db
db:
	docker compose up -d db

.PHONY: connect-db
connect-db:
	docker compose exec db psql -U omigoto
