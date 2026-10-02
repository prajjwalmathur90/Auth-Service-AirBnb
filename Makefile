include .env


dev:
	air

run:
	go run .

build:
	go build -o bin/app .



# Create new migration ---> make migrate-create name="new_migration_name"
migrate-create: 
	goose -dir $(MIGRATION_FOLDER) create $(name) sql


migrate-up:
	goose -dir $(MIGRATION_FOLDER) mysql $(DB_URL) up


migrate-down:
	goose -dir $(MIGRATION_FOLDER) mysql $(DB_URL) down


migrate-redo:
	goose -dir $(MIGRATION_FOLDER) mysql $(DB_URL) redo


migrate-reset:
	goose -dir $(MIGRATION_FOLDER) mysql $(DB_URL) reset


migrate-status:
	goose -dir $(MIGRATION_FOLDER) mysql $(DB_URL) status


migrate-help:
	goose -h


migrate-to:
	goose -dir $(MIGRATION_FOLDER) mysql $(DB_URL) up-to $(version)


migrate-down-to:
	goose -dir $(MIGRATION_FOLDER) mysql $(DB_URL) down-to $(version)

	
migrate-force:
	goose -dir $(MIGRATION_FOLDER) mysql $(DB_URL) force $(version)