.DEFAULT_GOAL := help

GO := go

.PHONY: help
help:
	@echo "Comandos disponibles:"
	@echo "  make run     - Arranca la API"
	@echo "  make build   - Compila el binario en ./bin/blog-api"
	@echo "  make test    - Ejecuta tests"
	@echo "  make lint    - Ejecuta golangci-lint"
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

.PHONY: lint
lint:
	./scripts/lint.sh

.PHONY: tidy
tidy:
	$(GO) mod tidy
