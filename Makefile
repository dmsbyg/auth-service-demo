include .env
all: build

build:
	@echo "Building..."
	@go build -o main cmd/api/main.go

run:
	@go run cmd/api/main.go

test:
	@echo "Testing..."
	@go test ./... -v

clean:
	@echo "Cleaning..."
	@rm -f main

watch:
	@if command -v air > /dev/null; then \
	    air; \
	    echo "Watching...";\
	else \
	    read -p "Go's 'air' is not installed on your machine. Do you want to install it? [Y/n] " choice; \
	    if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
	        go install github.com/cosmtrek/air@latest; \
	        air; \
	        echo "Watching...";\
	    else \
	        echo "You chose not to install air. Exiting..."; \
	        exit 1; \
	    fi; \
	fi

new_migration:
	@if command -v migrate > /dev/null; then \
		echo "migrating"; \
		migrate create -ext sql -dir ./database/migrations -seq $@; \
	else \
		echo "golang-migrate is not installed on your machine. You can install it referring to this instruction: https://github.com/golang-migrate/migrate/blob/master/cmd/migrate/README.md"; \
	fi

migrateup:
	@echo "Migrating..."
	@go run cmd/migrate/main.go

drop_db:
	@echo "Dropping database..."

	@rm $(DB_URL)

docker_build:
	@docker build -t auth-service-demo .

docker_run:
	@echo "Running app from docker"
	@docker run -p 8081:8081 --rm -it auth-service-demo

.PHONY: all build run test clean new_migration migrateup drop_db docker_build docker_run

