package main

import (
	"bytes"
	"embed"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
)

//go:embed templates/*
var templateFS embed.FS

type PageData struct {
	CSS         template.CSS
	Heading     string
	Description string
	URL         string
}

var tmpl = sync.OnceValue(func() *template.Template {
	t, err := template.ParseFS(templateFS, "templates/template.html")
	if err != nil {
		log.Fatalf("Failed to parse template: %v", err)
	}
	return t
})

var css = sync.OnceValue(func() template.CSS {
	data, err := templateFS.ReadFile("templates/style.css")
	if err != nil {
		log.Fatalf("Failed to read CSS: %v", err)
	}
	return template.CSS(data)
})

var baseURL = sync.OnceValue(func() string {
	if os.Getenv("NAIS_CLUSTER_NAME") == "prod-gcp" {
		return "https://mine-klager.nav.no"
	}
	return "https://mine-klager.ansatt.dev.nav.no"
})

var nbHTML = renderToBytes(PageData{
	CSS:         css(),
	Heading:     "Mine saker hos Klageinstans",
	Description: "Her kan du se status på dine saker hos Klageinstans.",
	URL:         baseURL(),
})

var nnHTML = renderToBytes(PageData{
	CSS:         css(),
	Heading:     "Mine saker hjå Klageinstans",
	Description: "Her kan du sjå status på dine saker hjå Klageinstans.",
	URL:         baseURL() + "/nn",
})

var enHTML = renderToBytes(PageData{
	CSS:         css(),
	Heading:     "My cases with Nav Complaints Unit",
	Description: "Here you can see the status of your cases with Nav Complaints Unit (Klageinstans).",
	URL:         baseURL() + "/en",
})

func renderToBytes(data PageData) []byte {
	var buf bytes.Buffer
	if err := tmpl().Execute(&buf, data); err != nil {
		log.Fatalf("Failed to render template: %v", err)
	}
	return buf.Bytes()
}

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
	mux.HandleFunc("/fallback", serveHTML(nbHTML))
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
