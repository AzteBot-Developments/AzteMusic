# LOCAL DEVELOPMENT UTILITY SHELL APPS
up:
	docker compose up -d --remove-orphans --build

down:
	docker compose down -v

update-lavalink:
	openssl base64 -A -in internal/lavalink/application.prod.yml -out internal/lavalink/application.prod.yml.out