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

// pingMessage is the message that we send to clients.
type pingMessage struct {
	Message string `json:"message"`
	Counter int    `json:"counter"`
}

// handler is the HTTP handler that sends incremental updates to the client.
// It sends a "ping" message every second.
func handler(w http.ResponseWriter, r *http.Request) {
	sw := sse.NewResponseWriter(w, sse.Options{})
	fmt.Println("New client connected")

	i := 0

	for {
		i++
		select {
		case <-r.Context().Done():
			fmt.Println("A client has disconnected.")
			return
		case <-time.After(1 * time.Second):
			err := sw.Write("ping", pingMessage{"pong", i})
			if err != nil {
				fmt.Println("Error writing to client:", err)
				return
			}
		}
	}
}
