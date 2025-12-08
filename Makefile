.PHONY: run brun app-logs lint

run:
	@docker compose up -d

brun:
	@docker compose up --build -d

app-logs:
	@docker logs subs-app

lint:
	@golangci-lint run
