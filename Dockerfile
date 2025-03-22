FROM golang:1.24.1-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o url-shortener

ENV DB_PATH=/app/data.db

EXPOSE 8080

CMD ["./url-shortener"]
