package main

func handler(w http.ResponseWriter, r *http.request) {
	fmt.Fprintf(w, "Hello", r.URL.Path[1:])
}

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
