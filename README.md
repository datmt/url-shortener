# URL Shortener (Go + SQLite + Swagger UI)

A simple self-hosted URL shortener written in Go with:

- ✅ SQLite persistence
- ✅ Basic Auth protected CRUD APIs
- ✅ Admin key–based user management
- ✅ Redirection via `/r/{handle}`

---

## 🏁 Quick Start

### 1. Clone and Build
```bash
git clone https://github.com/datmt/url-shortener.git
cd url-shortener
go mod tidy
go build -o url-shortener
```

### 2. Set Environment Variables
```bash
export ADMIN_KEY=your_admin_key
export DB_PATH=./data.db
export PORT=8080
```

### 3. Run the Server
```bash
./url-shortener
```
Server will start on [http://localhost:8080](http://localhost:8080)

---

## 🔐 Authentication

### Basic Auth for Short Link APIs
Users must be created via the admin endpoint and authenticate using HTTP Basic Auth.

### Admin Key for User Creation
To create new users:
```bash
POST /admin/create-user
Header: X-Admin-Key: your_admin_key
Body: {"username": "alice", "password": "secret"}
```

---

## 🔗 Endpoints

| Method | Path                 | Auth      | Description                     |
|--------|----------------------|-----------|---------------------------------|
| POST   | `/shorten`           | Basic Auth | Create or update short link     |
| GET    | `/shorten/{handle}` | Basic Auth | Get target URL                  |
| DELETE | `/delete/{handle}`  | Basic Auth | Delete short link (owner only)  |
| GET    | `/r/{handle}`        | None       | Redirect to target              |
| POST   | `/admin/create-user` | Admin Key | Create new user                 |
| GET    | `/swagger.yaml`      | None       | Raw OpenAPI YAML                |
| GET    | `/docs/`             | None       | Interactive Swagger UI          |

---

## 🧪 API Testing

Use `curl`, Postman

```bash
curl -X POST http://localhost:8080/admin/create-user \
  -H "X-Admin-Key: your_admin_key" \
  -H "Content-Type: application/json" \
  -d '{"username": "alice", "password": "mypassword"}'
```


## 🐳 Docker

```Dockerfile
FROM golang:1.21-alpine
WORKDIR /app
COPY . .
RUN go mod tidy && go build -o url-shortener
ENV DB_PATH=/app/data.db
EXPOSE 8080
CMD ["./url-shortener"]
```

Build and run:
```bash
docker build -t url-shortener .
docker run -p 8080:8080 \
  -e ADMIN_KEY=your_admin_key \
  -e DB_PATH=/app/data.db \
  url-shortener
```

---
<!-- Link to datmt.com -->
Created by datmt at [datmt.com](https://datmt.com)
