FROM golang:1.25-alpine AS builder
#install a c compiler for cgo
RUN apk add --no-cache gcc musl-dev

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -o boredgames-bot .
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/boredgames-bot .
COPY --from=builder /app/.env .

RUN mkdir -p /data

CMD ["./boredgames-bot"]