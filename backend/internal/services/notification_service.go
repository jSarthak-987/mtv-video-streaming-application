package service

import (
	"fmt"
	"log"
	"net/http"
	"sync"
)

// clients is a map that keeps track of all active clients connected via Server-Sent Events (SSE).
// The keys are channels that send status updates, and the values are empty structs (used for set-like behavior).
var clients = make(map[chan string]struct{})

// clientsMu is a mutex used to protect concurrent access to the clients map.
var clientsMu sync.Mutex

// currClientChan is a variable holding the current client channel.
// This variable is updated every time a new client connects.
var currClientChan chan string

// GetCurrClientChan returns the current client channel.
// This function allows access to the latest channel being used for communication.
func GetCurrClientChan() chan string {
	return currClientChan
}

// AddClient adds a new client channel to the clients map.
// It locks the map to ensure thread-safe access while adding the client channel.
func AddClient(clientChan chan string) {
	clientsMu.Lock()
	clients[clientChan] = struct{}{}
	clientsMu.Unlock()
}

// RemoveClient removes a client channel from the clients map and closes it.
// It locks the map to ensure thread-safe access while removing the client channel.
func RemoveClient(clientChan chan string) {
	clientsMu.Lock()
	delete(clients, clientChan)
	clientsMu.Unlock()
}

// SendStatusUpdateToClient sends a status update to a specific client channel.
// If the client channel no longer exists, it logs a message and does nothing.
// If the channel is ready to receive, it sends the status; otherwise, it logs that the client is not ready.
func SendStatusUpdateToClient(clientChan chan string, status string) {
	clientsMu.Lock()
	defer clientsMu.Unlock()

	// Check if the client channel still exists in the map
	if _, ok := clients[clientChan]; !ok {
		log.Printf("Attempted to send status '%s' to a closed channel\n", status)
		return
	}

	// Send the status update to the specified client channel
	select {
	case clientChan <- status:
		log.Printf("Sent: '%s'\n", status)
	default:
		log.Printf("Client is not ready to receive the status '%s'\n", status)
	}
}

// StatusStreamHandler handles incoming SSE connections for status updates.
// It sets up HTTP headers for SSE, registers the client channel, and listens for updates.
// When a client disconnects, the channel is removed and closed.
func StatusStreamHandler(w http.ResponseWriter, r *http.Request) chan string {
	// Set up HTTP headers required for SSE
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Create a new channel for the client
	clientChan := make(chan string, 100)
	currClientChan = clientChan

	// Register the client channel in the clients map
	AddClient(clientChan)
	defer func() {
		// Unregister the client channel and close it upon disconnection
		RemoveClient(clientChan)
		close(clientChan)
	}()

	for {
		select {
		case msg, ok := <-clientChan:
			// Check if the channel is still open
			if !ok {
				log.Println("Client channel closed")
				return clientChan
			}
			// Send the message to the client via SSE
			fmt.Fprintf(w, "data: %s\n\n", msg)
			w.(http.Flusher).Flush() // Flush the response to ensure it's sent immediately
		case <-r.Context().Done():
			// Handle client disconnection
			log.Println("Client disconnected")
			return clientChan
		}
	}
}
