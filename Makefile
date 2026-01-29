.PHONY: run

# Variables
APP_NAME=employee-api
MAIN_PATH=./cmd/server

# Go commands
GO=go

# Run the application locally
run:
	$(GO) run $(MAIN_PATH)/main.go
