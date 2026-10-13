# SnippetBox

A small Go web app for sharing code snippets. Users can sign up, log in, and create snippets that expire after a set number of days. The home page lists the ten most recent snippets that have not expired.

## Stack

- Go (`net/http`, `html/template`, `database/sql`)
- MySQL
- [scs](https://github.com/alexedwards/scs) sessions stored in MySQL
- bcrypt password hashes
- CSRF protection with [nosurf](https://github.com/justinas/nosurf)
- HTTPS with a local TLS certificate

## Routes

| Method | Path | Auth |
| --- | --- | --- |
| GET | `/` | public |
| GET | `/snippet/view/{id}` | public |
| GET | `/user/signup` | public |
| POST | `/user/signup` | public |
| GET | `/user/login` | public |
| POST | `/user/login` | public |
| GET | `/snippet/create` | signed in |
| POST | `/snippet/create` | signed in |
| POST | `/user/logout` | signed in |

## Prerequisites

- Go 1.26 or later
- Docker (for MySQL) or a local MySQL server
- A TLS certificate pair at `tls/cert.pem` and `tls/key.pem`

## Database

`docker-compose.yml` starts MySQL 8 on port 3306 with root password `my-secret-pw` and a database named `mydb`. The app’s default DSN expects a database named `snippetbox` and a user `web` with password `pass`. Create that user and the tables after the container is up:

```bash
docker compose up -d
```

```sql
CREATE DATABASE snippetbox CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

CREATE USER 'web'@'%' IDENTIFIED BY 'pass';
GRANT SELECT, INSERT, UPDATE, DELETE ON snippetbox.* TO 'web'@'%';
FLUSH PRIVILEGES;

USE snippetbox;

CREATE TABLE snippets (
    id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
    title VARCHAR(100) NOT NULL,
    content TEXT NOT NULL,
    created DATETIME NOT NULL,
    expires DATETIME NOT NULL
);

CREATE INDEX idx_snippets_created ON snippets(created);

CREATE TABLE users (
    id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    hashed_password CHAR(60) NOT NULL,
    created DATETIME NOT NULL,
    CONSTRAINT users_uc_email UNIQUE (email)
);

CREATE TABLE sessions (
    token CHAR(43) PRIMARY KEY,
    data BLOB NOT NULL,
    expiry TIMESTAMP(6) NOT NULL
);

CREATE INDEX sessions_expiry_idx ON sessions (expiry);
```

Connect as root with:

```bash
docker exec -it mysql-db mysql -uroot -pmy-secret-pw
```

## TLS certificate

The server listens with TLS. For local development, generate a self-signed certificate:

```bash
mkdir -p tls
openssl req -x509 -newkey rsa:2048 -keyout tls/key.pem -out tls/cert.pem -days 365 -nodes -subj "/CN=localhost"
```

Your browser will warn about the self-signed certificate. Accept it to continue, or use [mkcert](https://github.com/FiloSottile/mkcert) if you want a locally trusted cert.

## Run

```bash
go run ./cmd/web
```

Flags:

| Flag | Default | Purpose |
| --- | --- | --- |
| `-addr` | `:4000` | HTTP address |
| `-dsn` | `web:pass@/snippetbox?parseTime=true` | MySQL data source name |

Open [https://localhost:4000](https://localhost:4000).

If MySQL is running in Docker and the client connects over TCP, pass a host in the DSN:

```bash
go run ./cmd/web -dsn "web:pass@tcp(127.0.0.1:3306)/snippetbox?parseTime=true"
```

## Project layout

```
cmd/web/          HTTP server, handlers, middleware, templates
internal/models/  Snippet and user database access
internal/validator/
ui/html/          HTML templates
ui/static/        CSS and JavaScript
```
