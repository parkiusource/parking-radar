package hub

import (
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

// WebSocketHub handles broadcasting messages to connected clients
type WebSocketHub struct {
	clients   sync.Map         // Concurrent map for clients
	broadcast chan interface{} // Channel for broadcasting messages
	stop      chan struct{}    // Channel to stop the hub
	mutex     sync.Mutex       // Mutex for critical sections (if needed)
}

// NewWebSocketHub initializes a new WebSocketHub
func NewWebSocketHub() *WebSocketHub {
	return &WebSocketHub{
		broadcast: make(chan interface{}, 20),
		stop:      make(chan struct{}),
	}
}

// Run listens for broadcast messages or a stop signal
func (hub *WebSocketHub) Run() {
	for {
		select {
		case msg := <-hub.broadcast:
			hub.broadcastMessage(msg)

		case <-hub.stop:
			log.Println("Stopping WebSocket hub")
			return
		}
	}
}

// broadcastMessage sends a message to all connected clients
func (hub *WebSocketHub) broadcastMessage(msg interface{}) {
	hub.clients.Range(func(key, value interface{}) bool {
		client := key.(*websocket.Conn)
		go func(c *websocket.Conn) {
			if err := c.WriteJSON(msg); err != nil {
				log.Println("Error sending message:", err)
				err := c.Close()
				if err != nil {
					return
				}
				hub.RemoveClient(c)
			}
		}(client)
		return true
	})
}

// Stop sends a signal to stop the hub
func (hub *WebSocketHub) Stop() {
	close(hub.stop)
}

// Broadcast sends a message to all connected clients
func (hub *WebSocketHub) Broadcast(message interface{}) {
	select {
	case hub.broadcast <- message:
	default:
		log.Println("Broadcast channel full, message dropped")
	}
}

// AddClient adds a new client to the hub
func (hub *WebSocketHub) AddClient(conn *websocket.Conn) {
	hub.clients.Store(conn, true)
	log.Printf("Client added. Total clients: %d\n", hub.countClients())
}

// RemoveClient removes a client from the hub
func (hub *WebSocketHub) RemoveClient(conn *websocket.Conn) {
	hub.clients.Delete(conn)
	log.Printf("Client removed. Total clients: %d\n", hub.countClients())
}

// countClients returns the number of connected clients
func (hub *WebSocketHub) countClients() int {
	count := 0
	hub.clients.Range(func(_, _ interface{}) bool {
		count++
		return true
	})
	return count
}
