# LOCAL DEVELOPMENT UTILITY SHELL APPS
up:
	docker compose up -d --remove-orphans --build

down:
	docker compose down -v