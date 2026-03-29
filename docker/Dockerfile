FROM golang:1.23-alpine AS builder
WORKDIR /app

# go.sum  ->
COPY go.mod ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o skills-service ./cmd/server/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /app/skills-service .
COPY --from=builder /app/config ./config
COPY --from=builder /app/internal/infra/postgres/migrations ./internal/infra/postgres/migrations
EXPOSE 8082
CMD ["./skills-service"]
