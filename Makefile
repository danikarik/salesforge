install-tools:
	@echo "Installing tools..."
	go install github.com/pressly/goose/v3/cmd/goose@latest

dev:
	@echo "Running development server..."
	go run cmd/api/main.go

create-db:
	@echo "Creating database..."
	createdb -U postgres salesforge_development
	createdb -U postgres salesforge_test

drop-db:
	@echo "Dropping database..."
	dropdb -U postgres salesforge_development
	dropdb -U postgres salesforge_test

migrate:
	@echo "Applying test migrations..."
	goose up
	GOOSE_DBSTRING=$(TEST_DATABASE_URL) goose up
