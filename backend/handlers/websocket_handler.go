package handlers

import (
	"log"
	"net/http"
	"taboo-game/websocket"

	"github.com/gin-gonic/gin"
	gorilla "github.com/gorilla/websocket"
)

var upgrader = gorilla.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // In production, you should check the origin
	},
}

type WebSocketHandler struct {
	manager *websocket.Manager
}

func NewWebSocketHandler(manager *websocket.Manager) *WebSocketHandler {
	return &WebSocketHandler{
		manager: manager,
	}
}

func (h *WebSocketHandler) HandleConnection(c *gin.Context) {
	gameID := c.Param("gameId")
	clientID := c.Query("clientId")

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}

	client := &websocket.Client{
		ID:      clientID,
		GameID:  gameID,
		Socket:  conn,
		Send:    make(chan []byte, 256),
		Manager: h.manager,
	}

	h.manager.Register(client)

	// Start goroutines for reading and writing
	go client.Read()
	go client.Write()
}
