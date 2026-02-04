FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o boredgames-bot .
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/boredgames-bot .
COPY --from=builder /app/.env .
CMD ["./boredgames-bot"]