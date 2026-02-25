package main

import (
	"crypto/tls"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/newrelic/go-agent/v3/newrelic"
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
	mu    sync.Mutex
	logs  []OrderLog
	nrApp *newrelic.Application
)

// ---------- Utilities ----------

func randomString(n int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// ------------- NewRelic Function -----------------

func initNewRelic() {
	var err error

	nrApp, err = newrelic.NewApplication(
		newrelic.ConfigAppName("order-api"),
		newrelic.ConfigLicense(os.Getenv("NEW_RELIC_LICENSE_KEY")),
		newrelic.ConfigDistributedTracerEnabled(true),
	)

	if err != nil {
		log.Fatal("Failed to initialize New Relic: ", err)

	}
	log.Printf("New Relic initialized successfully")
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

// ---------- Background Log Generator ----------

func startLogGenerator() {
	envs := []string{"Dev", "PTest", "STest", "VTest", "Prod"}
	httpCodes := []int{200, 201, 400, 404, 500, 503}
	ports := []PortDef{
		{"HTTP", 80},
		{"HTTPS", 443},
		{"FTP", 21},
		{"SSH", 22},
	}

	go func() {
		for {
			code := httpCodes[rand.Intn(len(httpCodes))]
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

			if nrApp != nil {
				nrApp.RecordCustomMetric("Custom/Orders/Generated", 1)
				nrApp.RecordCustomMetric("Custom/Orders/HTTPCode/"+strconv.Itoa(code), 1)
			}
			time.Sleep(500 * time.Millisecond)
		}
	}()

}

// ---------- API Pull Handler (SUMO SAFE) ----------

func orderHandler(w http.ResponseWriter, r *http.Request) {

	var txn *newrelic.Transaction

	if nrApp != nil {
		txn = nrApp.StartTransaction("GET /order")
		defer txn.End()

		w = txn.SetWebResponse(w)
		txn.SetWebRequestHTTP(r)
	}

	// Always return 200 + JSON
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		json.NewEncoder(w).Encode([]OrderLog{})
		return
	}

	startStr := r.URL.Query().Get("start")
	endStr := r.URL.Query().Get("end")

	start, errStart := time.Parse(time.RFC3339, startStr)
	end, errEnd := time.Parse(time.RFC3339, endStr)

	// IMPORTANT: non-nil slice
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
		"API HIT from=%s method=%s start=%s end=%s count=%d",
		r.RemoteAddr, r.Method, startStr, endStr, len(result),
	)

	json.NewEncoder(w).Encode(result)
}

// ---------- Main ----------

func main() {
	rand.Seed(time.Now().UnixNano())

	initNewRelic()
	startLogGenerator()

	// Catch-all to avoid path mismatch issues
	http.HandleFunc("/order", orderHandler)
	http.HandleFunc("/", orderHandler)

	server := &http.Server{
		Addr: port,
		TLSConfig: &tls.Config{
			MinVersion: tls.VersionTLS10, // PoC only
			MaxVersion: tls.VersionTLS13,
			NextProtos: []string{
				"http/1.1",
			},
		},
	}

	log.Printf("Starting HTTPS server on %s", port)
	log.Fatal(server.ListenAndServeTLS(certFile, keyFile))
}
