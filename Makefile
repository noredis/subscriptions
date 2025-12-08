.PHONY: run brun app-logs lint stop down vdown db

run:
	@docker compose up -d

brun:
	@docker compose up --build -d

app-logs:
	@docker logs subs-app

lint:
	@golangci-lint run

stop:
	@docker compose stop

down:
	@docker compose down

vdown:
	@docker compose down -v

db:
	@docker exec -it subs-db psql -U postgres -d subs_db
