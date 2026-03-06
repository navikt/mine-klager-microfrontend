package main

import (
	"bytes"
	"embed"
	"log"
	"net/http"
	"os"
	"strconv"
)

//go:embed dist/*
var distFS embed.FS

func mustRead(name string) []byte {
	data, err := distFS.ReadFile(name)
	if err != nil {
		log.Fatalf("Failed to read %s: %v", name, err)
	}
	return data
}

func baseURL() string {
	if os.Getenv("NAIS_CLUSTER_NAME") == "prod-gcp" {
		return "https://mine-klager.nav.no"
	}
	return "https://mine-klager.ansatt.dev.nav.no"
}

func replaceBaseURL(html []byte) []byte {
	return bytes.ReplaceAll(html, []byte("{{BASE_URL}}"), []byte(baseURL()))
}

var nbHTML = replaceBaseURL(mustRead("dist/nb.html"))
var nnHTML = replaceBaseURL(mustRead("dist/nn.html"))
var enHTML = replaceBaseURL(mustRead("dist/en.html"))

func serveHTML(html []byte) http.HandlerFunc {
	contentLength := strconv.Itoa(len(html))

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Content-Length", contentLength)
		w.Write(html)
	}
}

var okBytes = []byte("OK")

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write(okBytes)
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/nb", serveHTML(nbHTML))
	mux.HandleFunc("/nn", serveHTML(nnHTML))
	mux.HandleFunc("/en", serveHTML(enHTML))
	mux.HandleFunc("/isAlive", healthHandler)
	mux.HandleFunc("/isReady", healthHandler)

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting server on port %s", port)
	log.Printf("Base URL: %s", baseURL())

	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
