# Order API – Polling API Simulation (HTTPS + Sumo HTTP Ingestion)

A lightweight Go HTTPS service that:

- Simulates an application generating structured logs
- Exposes a simple HTTP endpoint
- Pushes logs to Sumo Logic HTTP Source
- Simulates 10% error rate
- Runs securely over TLS 1.2+

This project is designed for:

- Polling API simulations
- Log ingestion testing
- Observability experiments
- Sumo Logic HTTP Source validation
- SRE portfolio demonstrations

---

## 🚀 Features

- Structured JSON log generation
- Simulated production traffic
- 10% random HTTP 500 error generation
- Direct push to Sumo HTTP Source
- HTTPS server (TLS 1.2 minimum)
- Lightweight PoC architecture
- Ready for polling tool integration

---

## 📦 Log Structure

Each request generates and pushes a structured log:

```json
{
  "timestamp": "2026-02-25T10:30:15Z",
  "service": "order-api",
  "env": "Prod",
  "level": "INFO",
  "orderId": "AbC123Xy",
  "httpCode": 200,
  "portName": "HTTPS",
  "port": 443,
  "message": "order processed"
}
```

---

## 🔎 Log Level Mapping

| HTTP Code | Level |
|-----------|--------|
| 500+      | ERROR  |
| 400–499   | WARN   |
| <400      | INFO   |

---

## 🔁 How It Works (Polling Model)

1. Client sends request to `/`
2. Service:
   - Generates structured log
   - Simulates occasional failure (10%)
   - Pushes log to Sumo HTTP Source
   - Returns JSON response
3. External system (or script) repeatedly polls endpoint
4. Logs are ingested into Sumo for analysis

This simulates a real application that generates logs while being polled by monitoring systems.

---

## 📤 Sumo Logic Integration

Logs are pushed to Sumo using HTTP Source.

### Sumo Endpoint

The endpoint is hardcoded in the application:

```
https://collectors.au.sumologic.com/receiver/v1/http/XXXXXXXX
```

### HTTP Request Details

- Method: `POST`
- Content-Type: `application/json`
- Timeout: 5 seconds

If Sumo returns a non-2xx status, it is logged to stdout.

---

## 🌐 API Endpoint

### `GET /`

Each request:

- Generates one log entry
- Pushes it to Sumo
- Returns the JSON log response

### Example

```bash
curl -k https://your-domain
```

Example response:

```json
{
  "timestamp": "2026-02-25T10:30:15Z",
  "service": "order-api",
  "env": "Prod",
  "level": "ERROR",
  "orderId": "XyZ987Ab",
  "httpCode": 500,
  "portName": "HTTPS",
  "port": 443,
  "message": "order processed"
}
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
certFile = "/etc/letsencrypt/live/order-api.fadikaba.com/fullchain.pem"
keyFile  = "/etc/letsencrypt/live/order-api.fadikaba.com/privkey.pem"
```

Valid TLS certificates are required.

---

## 🛠️ Running the Application

### 1️⃣ Initialize Module

```bash
go mod init order-api
```

### 2️⃣ Run Application

```bash
sudo go run main.go
```

If binding to port 443 without sudo:

```bash
sudo setcap 'cap_net_bind_service=+ep' ./order-api
```

---

## 🧪 Polling Script Example

Simulate monitoring system polling every second:

```bash
while true
do
  curl -sk https://localhost > /dev/null
  echo "Polled at $(date)"
  sleep 1
done
```

Or simulate burst polling:

```bash
for i in {1..100}
do
  curl -sk https://localhost > /dev/null
  echo "Request $i sent"
  sleep 0.5
done
```

---

## 📊 Observability Scenarios

Inside Sumo Logic you can:

- Track error rate (count where httpCode=500)
- Create metrics from logs
- Build dashboards
- Alert on error threshold
- Measure ingestion rate
- Analyze structured JSON fields

---

## 🎯 Use Cases

This project is ideal for:

- Sumo HTTP Source ingestion testing
- API polling simulation
- Synthetic traffic generation
- Error rate simulation
- TLS validation
- SRE lab experimentation
- Interview or portfolio demonstration

---

## ⚠️ Limitations

- Hardcoded Sumo endpoint
- No retry logic
- No batching
- No authentication
- No persistence
- Not production hardened
- Single instance only

---

## 📈 Future Improvements

- Move Sumo endpoint to environment variable
- Add retry with exponential backoff
- Add batching for performance
- Add correlation ID support
- Add OpenTelemetry tracing
- Add Prometheus metrics
- Add health endpoint
- Add Dockerfile

---

## 👤 Author

Fadi Kaba  
Senior SRE / Cloud Engineer  
Observability & Platform Engineering