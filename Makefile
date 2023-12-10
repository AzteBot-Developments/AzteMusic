# LOCAL DEVELOPMENT UTILITY SHELL APPS
up:
	docker compose up -d --remove-orphans --build

down:
	docker compose down -v

update-lavalink:
	openssl base64 -A -in internal/lavalink/application.prod.yml -out internal/lavalink/application.prod.yml.out

update-env:
	openssl base64 -A -in cmd/radio-service/.prod.env -out cmd/radio-service/.prod.env.out
	openssl base64 -A -in cmd/music-service/1/.prod.env -out cmd/music-service/1/.prod.env.out
	openssl base64 -A -in cmd/music-service/2/.prod.env -out cmd/music-service/2/.prod.env.out