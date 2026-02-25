# Order API – HTTPS Log Generator with New Relic Integration

A lightweight Go service that:

- Generates structured application logs in memory
- Exposes a time-filtered `/order` endpoint
- Runs over HTTPS (Let’s Encrypt TLS)
- Sends custom metrics to New Relic
- Designed for observability and ingestion testing (Sumo, New Relic, etc.)

This project is built as a Proof of Concept (PoC) for SRE and observability experimentation.

---

## 🚀 Features

- Structured JSON log generation
- Random environments (Dev, PTest, STest, VTest, Prod)
- Random HTTP codes (200, 201, 400, 404, 500, 503)
- Background log generator (every 500ms)
- Time-based log filtering via API
- HTTPS server on port 443
- New Relic custom metrics
- Thread-safe in-memory storage
- Memory guard to prevent unbounded growth

---

## 📦 Log Structure

Each generated log entry:

```json
{
  "timestamp": "2026-02-25T09:21:30Z",
  "service": "order-api",
  "env": "Prod",
  "level": "ERROR",
  "orderId": "a8K29xLp",
  "httpCode": 500,
  "portName": "HTTPS",
  "port": 443,
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

## 🔁 Background Log Generator

A goroutine runs continuously:

- Generates logs every 500ms
- Randomizes:
  - Environment
  - HTTP status code
  - Port (HTTP, HTTPS, FTP, SSH)
- Keeps a maximum of 5000 logs in memory
- Trims to last 3000 logs when threshold is exceeded
- Sends New Relic custom metrics:
  - `Custom/Orders/Generated`
  - `Custom/Orders/HTTPCode/<code>`

---

## 🌐 API Endpoint

### `GET /order`

Returns logs filtered by time range.

### Query Parameters

| Parameter | Required | Format |
|-----------|----------|--------|
| `start`   | Yes      | RFC3339 |
| `end`     | Yes      | RFC3339 |

### Example

```bash
curl -k "https://order-api.fadikaba.com/order?start=2026-02-25T09:00:00Z&end=2026-02-25T09:10:00Z"
```

### Behavior

- Always returns `200 OK`
- Returns empty JSON array if:
  - Invalid time format
  - Wrong HTTP method
  - No matching logs found

Example empty response:

```json
[]
```

---

## 🔐 HTTPS Configuration

The server runs on:

```
:443
```

Using Let’s Encrypt certificates:

```go
certFile = "/etc/letsencrypt/live/order-api.fadikaba.com/fullchain.pem"
keyFile  = "/etc/letsencrypt/live/order-api.fadikaba.com/privkey.pem"
```

TLS configuration:

- Minimum: TLS 1.0 (PoC only)
- Maximum: TLS 1.3
- HTTP/1.1 enabled

⚠️ For production use, set minimum TLS version to 1.2 or higher.

---

## 📊 New Relic Integration

The service initializes New Relic during startup.

### Required Environment Variable

```bash
export NEW_RELIC_LICENSE_KEY=YOUR_LICENSE_KEY
```

### Metrics Sent

- `Custom/Orders/Generated`
- `Custom/Orders/HTTPCode/200`
- `Custom/Orders/HTTPCode/400`
- `Custom/Orders/HTTPCode/500`
- etc.

Each `/order` request:

- Creates a New Relic transaction
- Adds custom attributes:
  - `clientIP`
  - `method`
  - `Path`

---

## 🧠 Concurrency Model

- Logs stored in a slice: `[]OrderLog`
- Access protected via `sync.Mutex`
- Background generator runs in a goroutine
- API safely reads logs under lock

Designed for PoC scale and observability testing.

---

## 🛠️ Running the Application

### 1. Initialize Module

```bash
go mod init order-api
go get github.com/newrelic/go-agent/v3/newrelic
```

### 2. Set Environment Variable

```bash
export NEW_RELIC_LICENSE_KEY=YOUR_KEY
```

### 3. Run

```bash
go build -o order-api
sudo setcap 'cap_net_bind_service=+ep' 
./order-api
```

---

## 🧪 Test Script Example

Send multiple test requests:

```bash
for i in {1..100}
do
  START=$(date -u -d "$i minutes ago" +"%Y-%m-%dT%H:%M:%SZ")
  END=$(date -u +"%Y-%m-%dT%H:%M:%SZ")

  curl -sk "https://localhost/order?start=$START&end=$END" > /dev/null
  echo "Request $i sent"
  sleep 0.5
done
```

---

## 🎯 Use Cases

This project is ideal for:

- Observability and ingestion testing
- New Relic custom metric experiments
- Sumo Logic Pull API testing
- TLS validation
- Load simulation
- SRE portfolio demonstrations
- Learning structured logging patterns

---

## ⚠️ Limitations

- In-memory storage only
- No persistence layer
- Not production hardened
- No horizontal scaling

---

## 📈 Future Improvements

- Add correlation ID support
- Add OpenTelemetry tracing
- Add Prometheus metrics endpoint
- Add health check endpoint
- Add Dockerfile
- Add authentication
- Improve TLS security configuration

---

## 👤 Author

Fadi Kaba  
Senior SRE / Cloud Engineer  
Observability & Platform Architecture  