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
	@echo "  make docker-build - Construye imagen Docker (Alpine)"
	@echo "  make docker-build-debian - Construye imagen Docker (Debian)"
	@echo "  make docker-build-all - Construye las dos imágenes Docker"
	@echo "  make docker-run - Levanta API + Prometheus"
	@echo "  make docker-stop - Detiene stack de Docker Compose"
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
	GOTOOLCHAIN=auto $(GO) run github.com/golangci/golangci-lint/cmd/golangci-lint@v1.64.8 run ./...

.PHONY: vuln
vuln:
	@which govulncheck >/dev/null 2>&1 || $(GO) install golang.org/x/vuln/cmd/govulncheck@latest
	govulncheck ./...

.PHONY: docker-build
docker-build:
	docker build -f deployments/Dockerfile -t gin-blog-api:latest .

.PHONY: docker-build-debian
docker-build-debian:
	docker build -f deployments/Dockerfile.debian -t gin-blog-api:debian .

.PHONY: docker-build-all
docker-build-all: docker-build docker-build-debian

.PHONY: docker-run
docker-run:
	docker compose -f deployments/docker-compose.yml up --build

.PHONY: docker-stop
docker-stop:
	docker compose -f deployments/docker-compose.yml down -v

.PHONY: tidy
tidy:
	$(GO) mod tidy
