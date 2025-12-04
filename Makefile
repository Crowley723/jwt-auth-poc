.PHONY: install dev debug test coverage coverage-html generate

install:
	go mod download

dev:
	GO_ENV=development reflex -r '\.go$$' -s -- go run ./main.go

debug:
	GO_ENV=development reflex -r '\.go$$' -s -- dlv debug --headless --listen=:2345 --api-version=2 --accept-multiclient ./main.go

TEST_FLAGS ?=

test:
	go test $(TEST_FLAGS) ./...

coverage:
	go test $(TEST_FLAGS) -cover ./...

coverage-html:
	go test -coverprofile=coverage.out $(TEST_FLAGS) ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

generate:
	go generate ./...

new-user:
	curl -X "POST" http://localhost:8080/api/users -H "Content-Type: application/json" -d '{"email":"test@example.com","name":"Test User"}'

delete-user:
	curl -X "DELETE" http://localhost:8080/api/users/1