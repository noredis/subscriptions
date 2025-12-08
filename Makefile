.PHONY: run brun app-logs lint stop down vdown

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
