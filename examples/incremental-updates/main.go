package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/floriscornel/sse-writer"
)

// main starts the server and listens for incoming connections.
func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe("localhost:8003", nil))
}

// handler is the HTTP handler that sends incremental updates to the client.
// It sends a list of players when a new client connects, and then sends
// random updates to the player list every 5 seconds.
func handler(w http.ResponseWriter, r *http.Request) {
	sw := sse.NewResponseWriter(w, sse.Options{
		Encoding: sse.Encode_Brotli,
	})
	fmt.Println("New client connected")

	// Send initial player list.
	select {
	case <-r.Context().Done():
		fmt.Println("A client has disconnected.")
	case <-time.After(1 * time.Second): // Mock delay.
		break
	}
	players := generatePlayerList(10)
	err := sw.Write(loadEvent, loadData{playerSliceFromMap(players)})
	if err != nil {
		fmt.Println("Error writing to client:", err)
		return
	}

	for {
		select {
		case <-r.Context().Done():
			fmt.Println("A client has disconnected.")
			return

		// Randomly add, update, or remove a player every 5 seconds.
		case <-time.After(5 * time.Second):
			err := performRandomEvent(sw, &players)
			if err != nil {
				fmt.Println("Error writing to client:", err)
				return
			}
		}
	}
}

// performRandomEvent performs a random event on the player list and sends the
// corresponding event to the client.
func performRandomEvent(sw sse.Writer, players *map[int]playerScore) error {
	events := []string{addEvent, updateEvent, removeEvent}
	if len(*players) == 0 {
		// If there are no players, only add a new player.
		events = []string{addEvent}
	}
	event := events[rand.Intn(len(events))]

	switch event {
	case addEvent:
		// Add a new player that doesn't already exist.
		for {
			newPlayer := generateRandomPlayer()
			if _, exists := (*players)[newPlayer.ID]; !exists {
				(*players)[newPlayer.ID] = newPlayer
				return sw.Write(addEvent, addData{newPlayer})
			}
		}
	case updateEvent:
		// Update a player's score.
		playerID := playerIDSliceFromMap(*players)[rand.Intn(len(*players))]
		player := (*players)[playerID]
		player.Score = generateRandomPlayerScore()
		(*players)[playerID] = player
		return sw.Write(updateEvent, updateData{player})
	case removeEvent:
		// Remove a player.
		playerID := playerIDSliceFromMap(*players)[rand.Intn(len(*players))]
		player := (*players)[playerID]
		delete(*players, playerID)
		return sw.Write(removeEvent, removeData{player})
	}
	return nil
}

const (
	// loadEvent is sent when a new client connects to the server.
	loadEvent = "load"
	// addEvent is sent when a new player is added.
	addEvent = "add"
	// updateEvent is sent when a player's score is updated.
	updateEvent = "update"
	// removeEvent is sent when a player is removed.
	removeEvent = "remove"
)

// playerScore represents a player's ID, name, and score.
type playerScore struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Score int    `json:"score"`
}

// loadData is the data attached to the loadEvent.
type loadData struct {
	Players []playerScore `json:"players"`
}

// addData is the data attached to the addEvent.
type addData = struct {
	Player playerScore `json:"player"`
}

// updateData is the data attached to the updateEvent.
type updateData = struct {
	Player playerScore `json:"player"`
}

// removeData is the data attached to the removeEvent.
type removeData = struct {
	Player playerScore `json:"player"`
}
