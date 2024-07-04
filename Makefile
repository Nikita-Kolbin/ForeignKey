compose_up:
	docker compose down -v
	docker compose up -d --build --force-recreate