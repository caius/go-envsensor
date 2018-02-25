package envsensor

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"text/template"
	"time"
)

type WebServer struct {
	listen        string
	cacheDuration time.Duration
	currentValue  CachedReading
}

func (h *WebServer) handleRoot(w http.ResponseWriter, r *http.Request) {
	log.Info("GET /")

	if h.currentValue.IsStale() {
		var s struct{}
		json.NewEncoder(w).Encode(s)
	} else {
		json.NewEncoder(w).Encode(h.currentValue.Reading)
	}
}

func (h *WebServer) handleMetrics(w http.ResponseWriter, r *http.Request) {
	log.Info("GET /metrics")

	const metricsTemplate = `
{{- if .Temperature -}}
# TYPE temperature gauge
temperature {{.Temperature}}
{{- end}}
{{ if .Humidity -}}
# TYPE humidity gauge
humidity {{.Humidity}}
{{- end}}
  `

	if h.currentValue.IsStale() {
		fmt.Fprintf(w, "\n")
	} else {
		tmpl, err := template.New("metrics").Parse(metricsTemplate)
		if err != nil {
			log.Fatal(err)
			fmt.Fprintf(w, "\n")
		}

		tmpl.Execute(w, h.currentValue.Reading)
	}
}

func (h *WebServer) subscribeToReadings(readings <-chan Reading) {
	log.Debug("WebServer subscribing to readings")
	for reading := range readings {
		log.WithFields(log.Fields{
			"reading": reading,
		}).Info("WebServer received reading")

		cachedReading := NewCachedReading(reading, h.cacheDuration)
		h.currentValue = cachedReading
	}
	log.Debug("WebServer finished listening for readings")
}

func NewWebServer(listen string, cacheDuration time.Duration) WebServer {
	return WebServer{
		listen:        listen,
		cacheDuration: cacheDuration,
	}
}

// Blocks and listens for http calls
func (h *WebServer) Start(readings <-chan Reading) {
	// Subscribe to readings & cache them for duration
	go h.subscribeToReadings(readings)

	// Handle HTTP calls
	http.HandleFunc("/", h.handleRoot)
	http.HandleFunc("/metrics", h.handleMetrics)

	log.Info("WebServer waiting to answer all your requests on ", h.listen)
	http.ListenAndServe(h.listen, nil)
}

// Noop
func (h *WebServer) Stop() {}
