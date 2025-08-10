DB_NAME=transfer_system_db
DB_USER=postgres
DB_PASS=postgres
DB_HOST=localhost
DB_PORT=5432

start-db:
	@docker start $(CONTAINER_NAME) 2>/dev/null || \
	docker run --name $(CONTAINER_NAME) \
		-e POSTGRES_USER=$(DB_USER) \
		-e POSTGRES_PASSWORD=$(DB_PASS) \
		-e POSTGRES_DB=$(DB_NAME) \
		-p $(DB_PORT):5432 -d $(POSTGRES_IMAGE)

# Stop Postgres container
stop-db:
	docker stop $(CONTAINER_NAME)

create-db:
	@echo "Creating database..."
	PGPASSWORD=$(DB_PASS) createdb -h $(DB_HOST) -p $(DB_PORT) -U $(DB_USER) $(DB_NAME) || true

drop-db:
	PGPASSWORD=$(DB_PASS) dropdb -h $(DB_HOST) -p $(DB_PORT) -U $(DB_USER) $(DB_NAME) || true

migrate:
	migrate -path ./migrations -database "postgres://$(DB_USER):$(DB_PASS)@localhost:5432/$(DB_NAME)?sslmode=disable" up

reset-db: stop-db start-db drop-db create-db migrate
