package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sync"

	"github.com/floriscornel/sse-writer"
)

// main starts the server and listens for incoming connections.
func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe("localhost:8002", nil))
}

// listenerMessage is the message that we send to listeners.
type listenerMessage struct {
	UserCount int `json:"user_count"`
}

// handler is the HTTP handler that sends incremental updates to the client.
// It sends a "update" message every time a new client connects or disconnects.
func handler(w http.ResponseWriter, r *http.Request) {
	sw := sse.NewResponseWriter(w, sse.Options{
		Encoding: sse.EncodeGzip,
	})

	// We generate a random ID for this listener.
	uniqueID := rand.Intn(1 << 31)
	fmt.Println("New listener connected with id", uniqueID)

	// Add the listener to the listeners map.
	addListener(uniqueID)
	defer func() {
		// Remove the listener from the listeners map when the connection is closed.
		removeListener(uniqueID)
	}()

	err := sw.Write("update", listenerMessage{len(listeners)})
	if err != nil {
		fmt.Println("Error writing to client:", err)
		return
	}

	for {
		select {
		case <-r.Context().Done():
			fmt.Println("The client has disconnected of listener", uniqueID)
			return
		case newCount := <-listeners[uniqueID]:
			err := sw.Write("update", listenerMessage{newCount})
			if err != nil {
				fmt.Println("Error writing to client:", err)
				return
			}
		}
	}
}

var (
	// listeners is a map of listener IDs to channels that we can write to.
	listeners = make(map[int]chan int)
	// listenerMutex is a mutex prevent concurrent access to the listeners map.
	listenerMutex sync.Mutex
)

// addListener adds a new listener to the listeners map and returns a channel that
func addListener(id int) chan int {
	listenerMutex.Lock()
	newCount := len(listeners) + 1
	for _, c := range listeners {
		c <- newCount
	}
	listeners[id] = make(chan int)
	listenerMutex.Unlock()
	return listeners[id]
}

// removeListener removes a listener from the listeners map.
func removeListener(deletionID int) {
	listenerMutex.Lock()
	newCount := len(listeners) - 1
	for id, c := range listeners {
		if id != deletionID {
			c <- newCount
		}
	}
	close(listeners[deletionID])
	delete(listeners, deletionID)
	listenerMutex.Unlock()
}
