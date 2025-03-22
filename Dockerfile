
FROM golang:1.24.1-alpine AS builder

ENV CGO_ENABLED=1
RUN apk add --no-cache gcc musl-dev sqlite

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .
RUN go build -o url-shortener

FROM alpine:3.18

RUN apk add --no-cache sqlite-libs ca-certificates

WORKDIR /app

COPY --from=builder /app/url-shortener /app/url-shortener
COPY --from=builder /app/ui /app/ui

ENV DB_PATH=/data/url-shortener.db
EXPOSE 8080
CMD ["./url-shortener"]
