package main

import (
        "bytes"
        "crypto/tls"
        "encoding/json"
        "log"
        "math/rand"
        "net/http"
        "os"
        "time"
)

const (
        port = ":443"

        certFile = "/home/ec2-user/development/cert/fullchain.pem"
        keyFile  = "/home/ec2-user/development/cert/privkey.pem"
)

var (
        sumoEndpoint = os.Getenv("SUMO_TOKEN")
)

// ---------- Log Structure ----------

type OrderLog struct {
        Timestamp     string `json:"timestamp"`
        Service       string `json:"service"`
        Env           string `json:"env"`
        Level         string `json:"level"`
        OrderID       string `json:"orderId"`
        HTTPCode      int    `json:"httpCode"`
        PortName      string `json:"portName"`
        Port          int    `json:"port"`
        CorrelationID string `json:"correlationId"`
        Message       string `json:"message"`
}

// ---------- Helpers ----------

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

// ---------- Push To Sumo ----------

func pushToSumo(logEntry OrderLog) error {

        payload, err := json.Marshal(logEntry)
        if err != nil {
                return err
        }

        req, err := http.NewRequest("POST", sumoEndpoint, bytes.NewBuffer(payload))
        if err != nil {
                return err
        }

        req.Header.Set("Content-Type", "text/plain")

        req.Header.Set("X-Sumo-Category", "order-api/prod")
        req.Header.Set("X-Sumo-Name", "direct-push")
        req.Header.Set("X-Sumo-Host", "Order-API")

        client := &http.Client{
                Timeout: 5 * time.Second,
        }

        resp, err := client.Do(req)
        if err != nil {
                return err
        }
        defer resp.Body.Close()

        if resp.StatusCode >= 300 {
                log.Printf("Sumo returned status: %s", resp.Status)
        }

        return nil
}

// Generate a random Correlation-ID

func generateCorrelationID(n int) string {

        genID := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
        b := make([]rune, n)

        for i := range b {
                b[i] = genID[rand.Intn(len(genID))]
        }

        return string(b)
}

// ---------- HTTP Handler ----------

func orderHandler(w http.ResponseWriter, r *http.Request) {

        correlationID := r.Header.Get("X-Correlation-ID")

        if correlationID == "" {
                correlationID = generateCorrelationID(12)
        } 

        // Simulate occasional error
        code := 200
        if rand.Float64() < 0.10 {
                code = 500
        }

        entry := OrderLog{
                Timestamp:     time.Now().UTC().Format(time.RFC3339),
                Service:       "order-api",
                Env:           "Prod",
                Level:         levelFromHTTPCode(code),
                OrderID:       randomString(8),
                CorrelationID: correlationID,
                HTTPCode:      code,
                PortName:      "HTTPS",
                Port:          443,
                Message:       "order processed",
        }
	log.Printf("CorrelationID being sent: %s", entry.CorrelationID)

        // Push to Sumo
        err := pushToSumo(entry)
        if err != nil {
                log.Printf("Push failed: %v", err)
        }

	w.Header().Set("X-Correlation-ID", correlationID)
        w.Header().Set("Content-Type", "application/json")

        json.NewEncoder(w).Encode(entry)
}

// ---------- Main ----------

func main() {
        rand.Seed(time.Now().UnixNano())

        http.HandleFunc("/", orderHandler)

        server := &http.Server{
                Addr: port,
                TLSConfig: &tls.Config{
                        MinVersion: tls.VersionTLS12,
                },
        }

        log.Printf("Starting HTTPS server on %s", port)
        log.Fatal(server.ListenAndServeTLS(certFile, keyFile))
}
