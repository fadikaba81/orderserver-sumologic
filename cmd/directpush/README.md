# Order API – Direct Push to Sumo Logic (HTTPS + Correlation ID)

A lightweight Go HTTPS service that:

- Accepts HTTP requests
- Generates structured JSON logs
- Injects or propagates a Correlation ID
- Pushes logs directly to Sumo Logic HTTP Source
- Simulates occasional 500 errors (10%)
- Runs over TLS

This project is designed for observability testing, correlation ID validation, and Sumo Logic ingestion experiments.

---

## 🚀 Features

- Structured JSON logging
- Direct push to Sumo Logic via HTTP Source
- Correlation ID propagation (X-Correlation-ID)
- Automatic Correlation ID generation if missing
- 10% simulated 500 errors
- HTTPS server (TLS 1.2 minimum)
- Custom Sumo metadata headers
- Lightweight and simple PoC architecture

---

## 📦 Log Structure

Each request generates and pushes the following JSON:

```json
{
  "timestamp": "2026-02-25T10:15:30Z",
  "service": "order-api",
  "env": "Prod",
  "level": "ERROR",
  "orderId": "a8K29xLp",
  "httpCode": 500,
  "portName": "HTTPS",
  "port": 443,
  "correlationId": "AbC123XyZ890",
  "message": "order processed"
}
```

### Log Level Mapping

| HTTP Code | Level |
|-----------|--------|
| 500+      | ERROR  |
| 400–499   | WARN   |
| <400      | INFO   |

---

## 🔗 Correlation ID Behavior

The service supports distributed tracing style correlation:

### Request Flow

1. If client sends:
   ```
   X-Correlation-ID: abc123
   ```
   → That value is used.

2. If header is missing:
   → A random 12-character ID is generated.

3. The Correlation ID:
   - Is included in the JSON log
   - Is returned in the response header
   - Is pushed to Sumo Logic

This allows end-to-end traceability in log platforms.

---

## 📤 Push to Sumo Logic

Logs are pushed directly to Sumo using an HTTP Source.

### Required Environment Variable

```bash
export SUMO_TOKEN=https://collectors.au.sumologic.com/receiver/v1/http/XXXXXXXX
```

The token must be a Sumo Logic HTTP Logs Source endpoint.

---

### Sumo Headers Used

```text
Content-Type: text/plain
X-Sumo-Category: order-api/prod
X-Sumo-Name: direct-push
X-Sumo-Host: Order-API
```

- `text/plain` is required for HTTP Logs Source
- Metadata headers allow categorization inside Sumo

---

## 🌐 API Endpoint

### `GET /`

Each request:

- Generates a structured log
- Simulates a 10% chance of HTTP 500
- Pushes the log to Sumo
- Returns the log as JSON response

### Example Request

```bash
curl -k https://your-domain
```

### Example With Correlation ID

```bash
curl -k https://your-domain \
  -H "X-Correlation-ID: my-custom-id"
```

### Response Example

```json
{
  "timestamp": "2026-02-25T10:15:30Z",
  "service": "order-api",
  "env": "Prod",
  "level": "INFO",
  "orderId": "ABcD1234",
  "httpCode": 200,
  "portName": "HTTPS",
  "port": 443,
  "correlationId": "my-custom-id",
  "message": "order processed"
}
```

Response header will also include:

```
X-Correlation-ID: my-custom-id
```

---

## 🔐 HTTPS Configuration

Server runs on:

```
:443
```

TLS Configuration:

- Minimum TLS version: 1.2
- Certificates loaded from:

```go
certFile = "/home/ec2-user/development/cert/fullchain.pem"
keyFile  = "/home/ec2-user/development/cert/privkey.pem"
```

You must provide valid TLS certificates.

---

## 🛠️ Running the Application

### 1️⃣ Initialize Go Module

```bash
go mod init order-api
```

### 2️⃣ Set Sumo Endpoint

```bash
export SUMO_TOKEN=https://collectors.au.sumologic.com/receiver/v1/http/YOUR_TOKEN
```

### 3️⃣ Run

```bash
sudo go run main.go
```

If binding to port 443 without sudo:

```bash
sudo setcap 'cap_net_bind_service=+ep' ./order-api
```

---

## 🧪 Load Testing Example

Generate 50 requests:

```bash
for i in {1..50}
do
  curl -sk https://localhost > /dev/null
  echo "Request $i sent"
  sleep 0.3
done
```

---

## 🎯 Use Cases

This project is ideal for:

- Sumo Logic HTTP Source testing
- Correlation ID validation
- Distributed tracing experiments
- Observability demonstrations
- Interview portfolio projects
- Error rate simulation testing
- Log ingestion validation

---

## 📊 Observability Notes

Inside Sumo Logic you can:

- Search by `_sourceCategory=order-api/prod`
- Filter by `correlationId`
- Build metrics from logs on `httpCode`
- Track error rate (500 responses)
- Create dashboards using correlationId

---

## ⚠️ Limitations

- No authentication
- No retry logic on failed push
- No circuit breaker
- No persistence
- Not production hardened
- No OpenTelemetry tracing
- Single instance only

---

## 📈 Future Improvements

- Add retry + exponential backoff for Sumo push
- Add batching for performance
- Add OpenTelemetry tracing
- Add health endpoint
- Add Prometheus metrics
- Add Dockerfile
- Add rate limiting
- Add structured logger (zap / zerolog)

---

## 👤 Author

Fadi Kaba  
Senior SRE / Cloud Engineer  
Observability & Platform Engineering
