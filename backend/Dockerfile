FROM golang:1.24-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

RUN go install github.com/swaggo/swag/cmd/swag@v1.16.6

COPY . .

RUN swag init -g cmd/server/main.go -o docs --parseInternal
RUN CGO_ENABLED=0 GOOS=linux go build -o /server ./cmd/server

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /server .
COPY --from=builder /app/migrations ./migrations

ENV MIGRATIONS_PATH=/app/migrations

EXPOSE 8080

CMD ["./server"]
