# syntax=docker/dockerfile:1.6

# ===== Build =====
FROM golang:1.23-alpine AS builder

RUN apk add --no-cache git
WORKDIR /build

COPY go.mod go.sum* ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux \
    go build -trimpath -ldflags="-s -w" -o /out/library-api ./cmd/server

# ===== Runtime =====
FROM alpine:3.20

# poppler-utils → pdftoppm (генерация обложек из PDF)
# wget оставляем для HEALTHCHECK
RUN apk add --no-cache poppler-utils ca-certificates tzdata wget \
 && addgroup -S library && adduser -S -G library -h /app library \
 && mkdir -p /var/library && chown -R library:library /var/library

WORKDIR /app
COPY --from=builder --chown=library:library /out/library-api /app/library-api
COPY --from=builder --chown=library:library /build/migrations /app/migrations

USER library

ENV API_PORT=2020
ENV MIGRATIONS_PATH=/app/migrations
ENV STORAGE_PATH=/var/library

EXPOSE 2020

HEALTHCHECK --interval=20s --timeout=5s --start-period=20s --retries=3 \
  CMD wget -qO- http://127.0.0.1:2020/api/health > /dev/null 2>&1 || exit 1

CMD ["/app/library-api"]
