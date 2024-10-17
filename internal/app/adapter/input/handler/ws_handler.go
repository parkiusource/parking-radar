package handler

import (
	"log"
	"net/http"

	"github.com/CamiloLeonP/parking-radar/internal/hub"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// WebSocketHandler encapsula la lógica de WebSocket y dependencias
type WebSocketHandler struct {
	hub *hub.WebSocketHub
}

// NewWebSocketHandler inicializa una nueva instancia de WebSocketHandler
func NewWebSocketHandler(hub *hub.WebSocketHub) *WebSocketHandler {
	return &WebSocketHandler{hub: hub}
}

// Upgrader convierte una conexión HTTP a WebSocket
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true }, // Permitir todas las conexiones
}

// HandleConnection maneja el ciclo de vida de la conexión WebSocket
func (wsh *WebSocketHandler) HandleConnection(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Failed to upgrade to WebSocket:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to upgrade to WebSocket"})
		return
	}
	defer func(conn *websocket.Conn) {
		if err := conn.Close(); err != nil {
			log.Println("Error closing WebSocket connection:", err)
		}
	}(conn)

	log.Println("New WebSocket connection established")

	// Enviar mensaje de bienvenida en formato JSON
	welcomeMessage := gin.H{
		"type":    "welcome",
		"payload": gin.H{"message": "Welcome to the WebSocket server for parking radar"},
	}
	if err := conn.WriteJSON(welcomeMessage); err != nil {
		log.Println("Error sending welcome message:", err)
		wsh.hub.RemoveClient(conn)
		return
	}

	// Agregar al cliente al hub
	wsh.hub.AddClient(conn)

	// Escuchar mensajes entrantes
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("WebSocket connection error:", err)
			wsh.hub.RemoveClient(conn)
			break
		}
		log.Printf("Received message: %s\n", message)
	}
}