package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

func mainHandler(w http.ResponseWriter, r *http.Request) {
	slog.Info("New session", slog.String("from", r.RemoteAddr))
	defer slog.Info("Session finished", slog.String("from", r.RemoteAddr))

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "SSE not supported", http.StatusInternalServerError)
		return
	}

	counter := 0
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if _, err := fmt.Fprintf(w, "%d\n", counter); err != nil {
				return
			}

			counter++
			flusher.Flush()
		case <-r.Context().Done():
			return
		}
	}
}

func main() {
	http.HandleFunc("/", mainHandler)
	if err := http.ListenAndServe(":80", nil); err != nil {
		panic(err)
	}
}
