.DEFAULT_GOAL := help

GO := go

.PHONY: help
help:
	@echo "Comandos disponibles:"
	@echo "  make run     - Arranca la API"
	@echo "  make build   - Compila el binario en ./bin/blog-api"
	@echo "  make test    - Ejecuta tests"
	@echo "  make coverage - Ejecuta tests con coverage y genera HTML"
	@echo "  make lint    - Ejecuta golangci-lint"
	@echo "  make vuln    - Ejecuta govulncheck"
	@echo "  make tidy    - Sincroniza dependencias del módulo"

.PHONY: run
run:
	$(GO) run ./cmd/api

.PHONY: build
build:
	@mkdir -p bin
	$(GO) build -o bin/blog-api ./cmd/api

.PHONY: test
test:
	$(GO) test ./...

.PHONY: coverage
coverage:
	$(GO) test -race -coverprofile=coverage.out -covermode=atomic ./...
	$(GO) tool cover -func=coverage.out | tail -n 1
	$(GO) tool cover -html=coverage.out -o coverage.html

.PHONY: lint
lint:
	./scripts/lint.sh

.PHONY: vuln
vuln:
	@which govulncheck >/dev/null 2>&1 || $(GO) install golang.org/x/vuln/cmd/govulncheck@latest
	govulncheck ./...

.PHONY: tidy
tidy:
	$(GO) mod tidy
