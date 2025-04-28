APP_NAME = myapp
BUILD_DIR = bin
GO_FILES = $(shell find be -type f -name '*.go' -not -path "./be/vendor/*")

.PHONY: all
all: build

.PHONY: build
build:
	@echo "🔨 Building backend..."
	cd be && go build -o ../$(BUILD_DIR)/$(APP_NAME) ./cmd/server

.PHONY: run
run:
	@echo "🚀 Running backend server..."
	cd be && clear && go run cmd/server/main.go

.PHONY: clean
clean:
	@echo "🧹 Cleaning build files..."
	rm -rf $(BUILD_DIR)

.PHONY: fmt
fmt:
	@echo "🧹 Formatting backend code..."
	cd be && go fmt ./...

.PHONY: test
test:
	@echo "🧪 Running backend tests..."
	cd be && go test ./...

.PHONY: tidy
tidy:
	@echo "🧹 Tidying modules..."
	cd be && go mod tidy

.PHONY: deps
deps:
	@echo "📦 Downloading backend dependencies..."
	cd be && go mod download