package envsensor

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"text/template"
	"time"
)

type WebServer struct {
	cacheDuration time.Duration
	currentValue  CachedReading
}

func (h *WebServer) handleRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Println("root request received")

	if h.currentValue.IsStale() {
		var s struct{}
		json.NewEncoder(w).Encode(s)
	} else {
		json.NewEncoder(w).Encode(h.currentValue.Reading)
	}
}

func (h *WebServer) handleMetrics(w http.ResponseWriter, r *http.Request) {
	fmt.Println("metric request received")

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
	for reading := range readings {
		fmt.Printf("Got reading! t=%f, h=%f\n", reading.Temperature, reading.Humidity)
		cachedReading := NewCachedReading(reading, h.cacheDuration)
		h.currentValue = cachedReading
	}
}

// Blocks and listens for http calls
func (h *WebServer) Start(readings <-chan Reading, cacheDuration time.Duration) {
	// Subscribe to readings & cache them for duration
	h.cacheDuration = cacheDuration
	go h.subscribeToReadings(readings)

	// Handle HTTP calls
	http.HandleFunc("/", h.handleRoot)
	http.HandleFunc("/metrics", h.handleMetrics)
	fmt.Printf("Waiting to answer all your requests on :8080\n")
	http.ListenAndServe(":8080", nil)
}
