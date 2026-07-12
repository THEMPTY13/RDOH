package main

import (
	"io"
	"log"
	"net/http"
)

const upstream = "https://cloudflare-dns.com/dns-query"

func dohHandler(w http.ResponseWriter, r *http.Request) {

	req, err := http.NewRequest(
		r.Method,
		upstream,
		r.Body,
	)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	req.Header.Set("Content-Type", "application/dns-message")
	req.Header.Set("Accept", "application/dns-message")

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		http.Error(w, err.Error(), 502)
		return
	}

	defer resp.Body.Close()

	w.Header().Set(
		"Content-Type",
		"application/dns-message",
	)

	w.WriteHeader(resp.StatusCode)

	io.Copy(w, resp.Body)
}

func main() {

	http.HandleFunc("/dns-query", dohHandler)

	log.Println("DoH running on :8080")

	log.Fatal(
		http.ListenAndServe(":8080", nil),
	)
}
