package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

var (
	mu          sync.Mutex
	logs        []string
	lastAddtime time.Time
)

type Ports struct {
	PortName   string
	PortNumber int
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func GenerateRandomString(n int) string {

	msg := "This is a logs from "
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	return msg + string(b)

}

func main() {

	env := []string{"Dev", "PTest", "STest", "VTest", "Prod"}
	httpCode := []int{200, 201, 400, 404, 500, 503}
	Port := []Ports{
		{
			PortName:   "HTTP",
			PortNumber: 80,
		},
		{
			PortName:   "HTTPS",
			PortNumber: 443,
		},
		{
			PortName:   "FTP",
			PortNumber: 21,
		},
		{
			PortName:   "SSH",
			PortNumber: 21,
		},
		{
			PortName:   "ICMP",
			PortNumber: -1,
		},
	}

	http.HandleFunc("/order", func(w http.ResponseWriter, r *http.Request) {

		t, err := time.LoadLocation("Australia/Melbourne")
		if err != nil {
			http.Error(w, "Invalid Timezone", http.StatusInternalServerError)
			return
		}

		now := time.Now()

		mu.Lock()

		if now.Sub(lastAddtime) > 500*time.Millisecond {
			x := rand.Intn((len(env) + 1) * 5)
			y := rand.Intn(len(Port))

			msg := GenerateRandomString(x)
			entry := fmt.Sprintf(
				"%s, %s,%d,%s,%d, %s, %d",
				now.In(t).Format(time.RFC3339),
				env[rand.Intn(len(env))],
				x,
				msg,
				httpCode[rand.Intn(len(httpCode))],
				Port[y].PortName,
				Port[y].PortNumber)

			logs = append(logs, entry)
			lastAddtime = now
		}

		for _, l := range logs {
			fmt.Fprintln(w, l)
		}

		mu.Unlock()

	})

	http.ListenAndServe(":8081", nil)
}
