package main

import (
        "encoding/json"
        "math/rand"
        "net/http"
        "sync"
        "time"
)

var (
        mu          sync.Mutex
        lastAddtime time.Time
        logs        []OrderLog
)

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

type Ports struct {
        PortName   string
        PortNumber int
}

func init() {
        rand.Seed(time.Now().UnixNano())
}

func GenerateRandomString(n int) string {
        letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
        b := make([]rune, n)
        for i := range b {
                b[i] = letters[rand.Intn(len(letters))]
        }
        return string(b)
}

func main() {

        envs := []string{"Dev", "PTest", "STest", "VTest", "Prod"}
        httpCodes := []int{200, 201, 400, 404, 500, 503}

        ports := []Ports{
                {"HTTP", 80},
                {"HTTPS", 443},
                {"FTP", 21},
                {"SSH", 22},
                {"ICMP", -1},
        }

        http.HandleFunc("/order", func(w http.ResponseWriter, r *http.Request) {

                loc, err := time.LoadLocation("Australia/Melbourne")
                if err != nil {
                        http.Error(w, "Invalid Timezone", http.StatusInternalServerError)
                        return
                }

                now := time.Now()

                mu.Lock()
                defer mu.Unlock()

                if now.Sub(lastAddtime) > 500*time.Millisecond {

                        env := envs[rand.Intn(len(envs))]
                        port := ports[rand.Intn(len(ports))]
                        httpCode := httpCodes[rand.Intn(len(httpCodes))]

                        entry := OrderLog{
                                Timestamp: now.In(loc).Format(time.RFC3339),
                                Service:   "order-api",
                                Env:       env,
                                Level:     levelFromHTTPCode(httpCode),
                                OrderID:   GenerateRandomString(8),
                                HTTPCode:  httpCode,
                                PortName:  port.PortName,
                                Port:      port.PortNumber,
                                Message:   "order processed",
                        }

                        logs = append(logs, entry)
                        lastAddtime = now
                }

                w.Header().Set("Content-Type", "application/json")
                json.NewEncoder(w).Encode(logs)
        })

        http.ListenAndServe(":8081", nil)
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
