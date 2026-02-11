package main

import (
        "crypto/tls"
        "encoding/json"
        "log"
        "math/rand"
        "net/http"
        "sync"
        "time"
)

const (
        port     = ":443"
        certFile = "/etc/letsencrypt/live/order-api.fadikaba.com/fullchain.pem"
        keyFile  = "/etc/letsencrypt/live/order-api.fadikaba.com/privkey.pem"
)

// ---------- Data Structures ----------

type OrderLog struct {
        Timestamp string `json:"timestamp"`
        Service   string `json:"service"`
        Env       string `json:"env"`
        Level     string `json:"level"`
        OrderID   string `json:"orderId"`
        HTTPCode  int    `json:"httpCode"`
        PortName  string `json:"portName"`
        Port      int    `json:"port"`
        Message   string `json:"message"`
}

type PortDef struct {
        Name   string
        Number int
}

// ---------- Globals ----------

var (
        mu   sync.Mutex
        logs []OrderLog
)

// Spike control
var (
        nextSpikeTime time.Time
        remaining5xx  int
)

// ---------- Helpers ----------

func random5xxDelay() time.Duration {
        min := 10 * time.Minute
        max := 20 * time.Minute
        return min + time.Duration(rand.Int63n(int64(max-min)))
}

func randomString(n int) string {
        letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
        b := make([]rune, n)
        for i := range b {
                b[i] = letters[rand.Intn(len(letters))]
        }
        return string(b)
}

func levelFromHTTPCode(code int) string {
        switch {
        case code >= 500:
                return "ERROR"
        case code >= 400:
                return "WARN"
        default:
                return "INFO"
        }
}

// ---------- Error Spike Logic ----------

func getCode() int {
        now := time.Now()

        // ---- Active spike: force 5xx ----
        if remaining5xx > 0 {
                remaining5xx--
                if rand.Intn(2) == 0 {
                        return 500
                }
                return 503
        }

        // ---- Start new spike ----
        if now.After(nextSpikeTime) && rand.Float64() < 0.02 {
                remaining5xx = 100 + rand.Intn(6) // 15–20 errors
                nextSpikeTime = now.Add(random5xxDelay())

                log.Printf("🔥 5xx SPIKE started (%d errors)", remaining5xx)

                // Emit first error immediately
                remaining5xx--
                if rand.Intn(2) == 0 {
                        return 500
                }
                return 503
        }

        // ---- Normal traffic ----
        if rand.Float64() < 0.80 {
                return 200
        }
        return 404
}

// ---------- Background Log Generator ----------

func startLogGenerator() {
        envs := []string{"Dev", "PTest", "STest", "VTest", "Prod"}
        ports := []PortDef{
                {"HTTP", 80},
                {"HTTPS", 443},
                {"FTP", 21},
                {"SSH", 22},
        }

        go func() {
                for {
                        code := getCode()
                        p := ports[rand.Intn(len(ports))]

                        entry := OrderLog{
                                Timestamp: time.Now().UTC().Format(time.RFC3339),
                                Service:   "order-api",
                                Env:       envs[rand.Intn(len(envs))],
                                Level:     levelFromHTTPCode(code),
                                OrderID:   randomString(8),
                                HTTPCode:  code,
                                PortName:  p.Name,
                                Port:      p.Number,
                                Message:   "order processed",
                        }

                        mu.Lock()
                        logs = append(logs, entry)

                        // PoC memory guard
                        if len(logs) > 5000 {
                                logs = logs[len(logs)-3000:]
                        }
                        mu.Unlock()

                        time.Sleep(500 * time.Millisecond)
                }
        }()
}

// ---------- API Handler ----------

func orderHandler(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")

        if r.Method != http.MethodGet {
                json.NewEncoder(w).Encode([]OrderLog{})
                return
        }

        startStr := r.URL.Query().Get("start")
        endStr := r.URL.Query().Get("end")

        start, errStart := time.Parse(time.RFC3339, startStr)
        end, errEnd := time.Parse(time.RFC3339, endStr)

        result := make([]OrderLog, 0)

        if errStart != nil || errEnd != nil {
                json.NewEncoder(w).Encode(result)
                return
        }

        mu.Lock()
        for _, l := range logs {
                ts, err := time.Parse(time.RFC3339, l.Timestamp)
                if err != nil {
                        continue
                }
                if !ts.Before(start) && ts.Before(end) {
                        result = append(result, l)
                }
        }
        mu.Unlock()

        log.Printf(
                "API HIT from=%s start=%s end=%s count=%d",
                r.RemoteAddr, startStr, endStr, len(result),
        )

        json.NewEncoder(w).Encode(result)
}

// ---------- Main ----------

func main() {
        rand.Seed(time.Now().UnixNano())

        startLogGenerator()

        http.HandleFunc("/order", orderHandler)
        http.HandleFunc("/", orderHandler)

        server := &http.Server{
                Addr: port,
                TLSConfig: &tls.Config{
                        MinVersion: tls.VersionTLS10, // PoC only
                        MaxVersion: tls.VersionTLS13,
                        NextProtos: []string{"http/1.1"},
                },
        }

        log.Printf("Starting HTTPS server on %s", port)
        log.Fatal(server.ListenAndServeTLS(certFile, keyFile))
}
