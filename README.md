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

Sample Output 

```
2026-01-29T21:45:31+11:00, Dev,21,This is a logs from 9MnVvFq25HbyHAOWDcX8W,200, SSH, 21
2026-01-29T21:45:33+11:00, PTest,6,This is a logs from XezOCQ,404, FTP, 21
2026-01-29T21:45:33+11:00, STest,16,This is a logs from N1ALVhNgU2X1pNrl,200, ICMP, -1
2026-01-29T21:45:34+11:00, Prod,6,This is a logs from coVbT7,201, HTTP, 80
2026-01-29T21:45:34+11:00, VTest,15,This is a logs from PIOCjm1TAddiknY,500, FTP, 21
2026-01-29T21:45:36+11:00, Dev,2,This is a logs from AM,400, HTTP, 80
2026-01-29T21:45:36+11:00, PTest,28,This is a logs from RqbuilNLufVvyAKne6gVQcSxrFti,400, HTTP, 80
2026-01-29T21:45:37+11:00, Prod,2,This is a logs from nQ,200, ICMP, -1
2026-01-29T21:45:38+11:00, Dev,10,This is a logs from LgxrAW8a4s,200, ICMP, -1
2026-01-29T21:45:38+11:00, VTest,12,This is a logs from xBdSwl6SzIsa,404, HTTP, 80
2026-01-29T21:45:39+11:00, STest,10,This is a logs from 4qGa1ihYQe,200, ICMP, -1
2026-01-29T21:45:39+11:00, Dev,21,This is a logs from zZAaonqOVsnT7D6sngrmR,400, ICMP, -1
2026-01-29T21:45:40+11:00, Prod,0,This is a logs from ,500, HTTPS, 443
2026-01-29T21:45:40+11:00, VTest,11,This is a logs from gAt7zq1JEfl,503, SSH, 21
2026-01-29T21:45:41+11:00, VTest,8,This is a logs from J97kRCAD,201, SSH, 21
2026-01-29T21:45:41+11:00, Prod,16,This is a logs from ui7Jh1G6N4KTib2Z,503, HTTP, 80
2026-01-29T21:45:42+11:00, PTest,4,This is a logs from Hxd6,200, ICMP, -1
2026-01-29T21:45:43+11:00, VTest,26,This is a logs from PaiejHIEaOQj7Nvv2wq8mpsCq0,500, HTTPS, 443
2026-01-29T21:45:43+11:00, Prod,1,This is a logs from R,200, HTTP, 80
2026-01-29T21:45:44+11:00, VTest,9,This is a logs from fZCr5SyOw,201, SSH, 21
2026-01-29T21:45:45+11:00, STest,27,This is a logs from gIJKfxkeuHfLuebuEkI7JZeSMdn,200, HTTP, 80
2026-01-29T21:45:45+11:00, Dev,17,This is a logs from R9wYi6HTLNbEDNLrj,404, ICMP, -1
2026-01-29T21:45:46+11:00, Prod,29,This is a logs from 2ZI2Hf6WsDxjJOXm7D5fe37lszACZ,400, HTTPS, 443
2026-01-29T21:45:46+11:00, VTest,27,This is a logs from w9ksLvssCclnopIkrZzD7wZ9dfV,404, SSH, 21
2026-01-29T21:45:47+11:00, PTest,2,This is a logs from kx,503, HTTPS, 443
2026-01-29T21:45:47+11:00, PTest,8,This is a logs from MPCtsjRZ,503, ICMP, -1
2026-01-29T21:45:48+11:00, STest,26,This is a logs from 594PevmzvEgY96vyoBNB9mfYgI,200, HTTPS, 443
2026-01-29T21:45:49+11:00, VTest,18,This is a logs from 3l3sDWZpH7pZcvCev7,201, ICMP, -1
2026-01-29T21:45:49+11:00, PTest,4,This is a logs from Luya,404, FTP, 21
2026-01-29T21:45:50+11:00, Prod,2,This is a logs from lf,404, HTTPS, 443
2026-01-29T21:45:50+11:00, PTest,22,This is a logs from 4elDY1F1IPDaTgkctK1PKJ,503, ICMP, -1
2026-01-29T21:45:51+11:00, PTest,7,This is a logs from q5hWKqq,201, ICMP, -1
2026-01-29T21:45:51+11:00, STest,7,This is a logs from Y9YvxvX,503, ICMP, -1
2026-01-29T21:45:52+11:00, VTest,3,This is a logs from tY9,200, HTTPS, 443
2026-01-29T21:45:52+11:00, STest,25,This is a logs from daMQAX10l21ymuYrkdwHTAXrk,503, HTTPS, 443
2026-01-29T21:45:53+11:00, VTest,7,This is a logs from 6KPBC2b,503, HTTP, 80
2026-01-29T21:45:53+11:00, STest,18,This is a logs from SMhYbisdC8QydPIrz8,503, HTTPS, 443
2026-01-29T21:45:54+11:00, STest,21,This is a logs from sP7nRYBBZeet83FpZaFNF,400, HTTPS, 443
2026-01-29T21:45:55+11:00, Dev,7,This is a logs from mHe0KqZ,404, SSH, 21
2026-01-29T21:45:55+11:00, Prod,25,This is a logs from jD8Rx5bR8jnPY7WwK83hsKrIY,500, SSH, 21
```