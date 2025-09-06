BINARY_NAME=avito-shop
PKG=./...
MAIN=./cmd/shop/main.go
COVER_PROFILE=coverage.out
COVER_HTML=coverage.html

.PHONY: all build run test clean tidy fmt lint tests_with_cover

all: build

build:
	go build -o $(BINARY_NAME) $(MAIN)

run: test build
	./$(BINARY_NAME)

test:
	go test -v $(PKG)

test_with_cover:
	go test -coverprofile=$(COVER_PROFILE) $(PKG)
	go tool cover -html=$(COVER_PROFILE) -o $(COVER_HTML)
	@echo "Coverage report generated: $(COVER_HTML)"

clean:
	rm -f $(BINARY_NAME) $(COVER_PROFILE) $(COVER_HTML)

fmt:
	go fmt $(PKG)

lint:
	go vet $(PKG)

tidy:
	go mod tidy
