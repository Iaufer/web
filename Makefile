.PHONY: build
build:
	go build -v ./cmd/apiserver


.PHONY: test
test:
	go test -v -timeout 30s ./...


.PHONY: migr_up
	@echo "Migrate up"
	migrate -path migrations -database "postgres://postgres:2505@localhost:5432/postgres?sslmode=disable" up

.PHONY: migr_down
	@echo "Migrate down"
	migrate -path migrations -database "postgres://postgres:2505@localhost:5432/postgres?sslmode=disable" down

.DEFAULT_GOAL := build