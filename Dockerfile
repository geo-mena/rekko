# 🚀 Stage 1: Build frontend
FROM node:20-alpine AS frontend-builder

WORKDIR /app

RUN corepack enable && corepack prepare pnpm@latest --activate

COPY frontend/package.json frontend/pnpm-lock.yaml ./
RUN pnpm install --frozen-lockfile

COPY frontend/ .

ENV VITE_SERVER_API_URL=""
ENV VITE_SERVER_API_PREFIX=/api/v1
ENV VITE_SERVER_API_TIMEOUT=180000

RUN pnpm build

# 🚀 Stage 2: Build backend
FROM golang:1.24-alpine AS backend-builder

WORKDIR /app

RUN apk add --no-cache git

COPY backend/go.mod backend/go.sum ./
RUN go mod download

RUN go install github.com/swaggo/swag/cmd/swag@v1.16.6

COPY backend/ .

RUN swag init -g cmd/server/main.go -o docs --parseInternal
RUN CGO_ENABLED=0 GOOS=linux go build -o /server ./cmd/server

# 🚀 Stage 3: Runtime
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=backend-builder /server .
COPY --from=backend-builder /app/migrations ./migrations
COPY --from=frontend-builder /app/dist ./static

ENV MIGRATIONS_PATH=/app/migrations
ENV STATIC_DIR=/app/static
ENV GIN_MODE=release

EXPOSE 8080

CMD ["./server"]
