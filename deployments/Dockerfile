FROM golang:1.25-alpine AS builder

RUN apk add --no-cache build-base sqlite-dev

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux GOMAXPROCS=2 \
    go build -v -p=1 -trimpath -ldflags="-s -w" -o /out/blog-api ./cmd/api

FROM alpine:3.20

LABEL MAINTAINER="Ivanruii"

RUN apk add --no-cache ca-certificates sqlite-libs tzdata wget

RUN addgroup -S app && adduser -S -G app app

WORKDIR /app

COPY --from=builder /out/blog-api /app/blog-api
RUN mkdir -p /app/data && chown -R app:app /app

USER app

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --spider -q http://localhost:8080/api/v1/health || exit 1

ENV GIN_MODE=release
ENV PORT=8080
ENV DB_PATH=/app/data/blog.db

VOLUME ["/app/data"]

ENTRYPOINT ["/app/blog-api"]
