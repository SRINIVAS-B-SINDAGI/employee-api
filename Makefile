.PHONY: all run test

# Variables
APP_NAME=employee-api
MAIN_PATH=./cmd/server
GOTEST=$(GO) test
MAIN_PATH=./cmd/server
GOBUILD=$(GO) build
GOMOD=$(GO) mod

# Go commands
GO=go

deps:
	$(GOMOD) download
	$(GOMOD) tidy

# Run the application locally
run:
	$(GO) run $(MAIN_PATH)/main.go

test:
	$(GOTEST) -v -race -cover ./...

test-coverage:
	$(GOTEST) -v -race -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html

gen-proto:
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		proto/auth/v1/*.proto \
		proto/employee/v1/*.proto \
		proto/salary/v1/*.proto

build:
	$(GOBUILD) -o bin/$(APP_NAME) $(MAIN_PATH)
	