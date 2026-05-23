# gin-blog-api

Proyecto de API REST para un blog utilizando el framework Gin en Go.

## Curls rápidos para probar funcionalidad

```bash
BASE_URL="http://localhost:18080"

# Health
curl -s "$BASE_URL/api/v1/health"

# Crear post
curl -s -X POST "$BASE_URL/api/v1/posts" \
  -H "Content-Type: application/json" \
  -d '{"title":"Mi primer post","content":"Contenido de prueba con más de diez chars","author":"Ivan","published":false}'

# Listar posts
curl -s "$BASE_URL/api/v1/posts?page=1&limit=10"

# Publicar post con ID 1
curl -s -X PATCH "$BASE_URL/api/v1/posts/1/publish"

# Crear comentario en post 1
curl -s -X POST "$BASE_URL/api/v1/posts/1/comments" \
  -H "Content-Type: application/json" \
  -d '{"author":"Ana","content":"Buen post!"}'

# Listar comentarios del post 1
curl -s "$BASE_URL/api/v1/posts/1/comments"
```

## Métricas Prometheus expuestas

Endpoint: `GET /metrics`

- `http_requests_total`
- `http_request_duration_seconds`
- `http_requests_in_flight`
- `http_errors_total`
- `posts_created_total`
- `posts_published_total`
- `posts_deleted_total`
- `comments_created_total`
- `posts_total`
- `comments_total`
- `db_query_duration_seconds`
- `db_errors_total`

Prometheus (contenedor `blog-prometheus`) scrapea este endpoint internamente contra `api:8080/metrics`.