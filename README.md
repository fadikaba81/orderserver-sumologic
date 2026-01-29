# Go Log Generator (Stdout → EC2 → Sumo Logic)

This project is a simple **Golang HTTP service** that generates **structured logs to stdout**, runs inside **EC2**, and is designed to be collected by the **Sumo Logic Universal (Installed) Collector**.

The goal is to demonstrate an **end‑to‑end logging pipeline** suitable for SRE / Observability use cases.

---

## 🧩 What This Project Does

* Runs a Go HTTP web server
* Exposes an endpoint (e.g. `/order`)
* On each request:

  * Generates a log entry containing:

    * timestamp
    * environment
    * HTTP status
    * message
  * Writes logs to **stdout** (not files)
* Runs on an Amazon EC2 instance
* Logs are collected externally by **Sumo Logic Universal Collector**

---

## 🏗 Architecture Overview

```
Client (Browser / curl)
        |
        v
Go HTTP Server (EC2 Instance)
        |
        v
Structured Logs (stdout)
        |
        v
Sumo Logic Universal Collector (Installed on EC2)
        |
        v
Sumo Logic Search / Dashboards
```

Client (Browser / curl)
|
v
EC2 Public IP
|
v
Structured Logs (stdout)
^
|
Sumo Logic Universal Collector
|
v
Sumo Logic Search / Dashboards

```

---

## 📁 Project Structure

```

go-log-generator/
├── cmd/main.go
├── go.mod
├── README.md


````

---

## ⚙️ Configuration

The application is configured using environment variables:

| Variable | Description | Example |
|--------|------------|--------|
| `APP_ENV` | Application environment | `dev`, `test`, `prod` |
| `APP_PORT` | HTTP port to listen on | `8081` |

If `APP_PORT` is not set, the app should default to a safe port (e.g. 8080 or 8081).

---

## 🚀 Running the Application Locally

### 1️⃣ Run with Go

```bash
go run main.go
````

### 2️⃣ Test the Endpoint

```bash
curl http://localhost:8081/order
```

### 3️⃣ Expected Results

* Browser / curl receives a plain-text or JSON response
* Terminal displays structured log output (stdout)

---

## 🖥 Running on an EC2 Instance

### 1️⃣ Provision EC2

* Launch an EC2 instance (Amazon Linux 2 or similar)
* Ensure inbound rules allow the application port (e.g. 8081)
* Install Go on the instance

### 2️⃣ Run the Application

```bash
go run main.go
```

Or build a binary:

```bash
go build -o go-log-generator
./go-log-generator
```

### 3️⃣ Test the Endpoint

```bash
curl http://<EC2_PUBLIC_IP>:8081/order
```

### 4️⃣ Verify Logs

Logs should appear in stdout and system logs on the EC2 instance.

---

## 🧾 Logging Design

* Logs are written to **stdout only**
* One log event per line
* No multiline logs
* Structured format (recommended: JSON)

### Example Log

```json
{
  "timestamp": "2026-01-29T15:22:10Z",
  "env": "dev",
  "http_status": 200,
  "message": "order processed successfully"
}
```

This format is intentionally designed for:

* Easy ingestion
* Field Extraction Rules (FER)
* Low parsing overhead in Sumo Logic

---

## 📡 Sumo Logic Collection Strategy

### Collector Type

* **Installed (Universal) Collector** running directly on the EC2 instance

### Log Source

The collector tails application logs written to stdout and captured via system logging (or redirected output if configured).

Each structured log line becomes a single Sumo Logic event.

---

## 🔎 Verifying in Sumo Logic

### Sample Search

```text
_sourceCategory=*go-log-generator*
```

### Example Queries

Count logs by environment:

```
| count by env
```

Count by HTTP status:

```
| count by http_status
```

---

## 🧠 Key Design Decisions

* **stdout logging** instead of file-based logging
* Collector installed directly on the EC2 instance
* Separation of:

  * User response
  * Application logs
  * Observability pipeline

This mirrors production-grade container logging best practices.

---

## 🚨 Common Pitfalls

* Using `log.Fatal` inside request handlers
* Writing logs to files instead of stdout
* Multi-line or unstructured logs
* Hardcoding ports or environments
* Not handling `ListenAndServe` errors

---

## 🎯 Future Enhancements

* Add structured JSON responses
* Add request IDs / correlation IDs
* Add OpenTelemetry tracing
* Deploy to Kubernetes
* Add dashboards and alerts in Sumo Logic

---

## 👤 Author

**Fadi Kaba**
Senior SRE / Cloud & Observability Engineer

---

## 📜 License

This project is provided for learning and demonstration purposes.
