# LOCAL DEVELOPMENT UTILITY SHELL APPS
up:
	docker compose up -d --remove-orphans --build

down:
	docker compose down -v

update-lavalink:
	openssl base64 -A -in internal/lavalink/application.prod.yml -out internal/lavalink/application.prod.yml.out

update-env:
	openssl base64 -A -in cmd/music-service/.prod.env -out cmd/music-service/.prod.env.out