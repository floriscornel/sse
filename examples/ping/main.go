package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/floriscornel/sse"
)

// main starts the server and listens for incoming connections.
func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe("localhost:8001", nil))
}

// handler is the HTTP handler that sends incremental updates to the client.
// It sends a "ping" message every second.
func handler(w http.ResponseWriter, r *http.Request) {
	opts := sse.Options{
		ResponseStatus: http.StatusOK,
		Encoding:       sse.EncodeNone,
	}
	sw := sse.NewResponseWriter(w, opts)

	for {
		select {
		case <-r.Context().Done():
			return
		case <-time.After(1 * time.Second):
			err := sw.Write("ping", "pong")
			if err != nil {
				fmt.Println("Error writing to client:", err)
				return
			}
		}
	}
}
