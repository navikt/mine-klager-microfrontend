package main

import (
	_ "embed"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

//go:embed templates/template.html
var templateHTML []byte

//go:embed templates/style.css
var styleCSS []byte

type variant struct {
	title       string
	description string
	urlSuffix   string
}

var variants = map[string]variant{
	"nb": {
		title:       "Mine saker hos Klageinstans",
		description: "Her kan du se status på dine saker hos Klageinstans.",
		urlSuffix:   "",
	},
	"nn": {
		title:       "Mine saker hjå Klageinstans",
		description: "Her kan du sjå status på dine saker hjå Klageinstans.",
		urlSuffix:   "/nn",
	},
	"en": {
		title:       "My cases with Nav Complaints Unit",
		description: "Here you can see the status of your cases with Nav Complaints Unit (Klageinstans).",
		urlSuffix:   "/en",
	},
}

func baseURL() string {
	if os.Getenv("NAIS_CLUSTER_NAME") == "prod-gcp" {
		return "https://mine-klager.nav.no"
	}
	return "https://mine-klager.ansatt.dev.nav.no"
}

func renderVariant(v variant) []byte {
	url := baseURL() + v.urlSuffix

	markup := string(templateHTML)
	markup = strings.ReplaceAll(markup, "{{TITLE}}", v.title)
	markup = strings.ReplaceAll(markup, "{{DESCRIPTION}}", v.description)
	markup = strings.ReplaceAll(markup, "{{URL}}", url)

	html := fmt.Sprintf(`
<mine-klager-microfrontend>
  <template shadowrootmode="open">
  <style>%s</style>
  %s
  </template>
</mine-klager-microfrontend>
<script>
  (function() {
    var el = document.currentScript.previousElementSibling;
    if (el.shadowRoot === null) {
      el.attachShadow({ mode: "open" }).appendChild(el.firstElementChild.content);
    }
  })()
</script>`, styleCSS, strings.TrimSpace(markup))

	return []byte(strings.TrimSpace(html))
}

var nbHTML = renderVariant(variants["nb"])
var nnHTML = renderVariant(variants["nn"])
var enHTML = renderVariant(variants["en"])

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
