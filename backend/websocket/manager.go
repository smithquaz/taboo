package websocket

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
	"taboo-game/types"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // In production, implement proper origin checking
	},
}

type Manager struct {
	mu              sync.RWMutex
	gameConnections map[string]map[*websocket.Conn]bool
	gameEvents      types.GameEventsServiceInterface
	register        chan types.WebSocketClientInterface
	unregister      chan types.WebSocketClientInterface
	shutdown        chan struct{}
}

func NewManager(gameEvents types.GameEventsServiceInterface) *Manager {
	return &Manager{
		gameConnections: make(map[string]map[*websocket.Conn]bool),
		gameEvents:      gameEvents,
		register:        make(chan types.WebSocketClientInterface),
		unregister:      make(chan types.WebSocketClientInterface),
		shutdown:        make(chan struct{}),
	}
}

func (m *Manager) Register(client types.WebSocketClientInterface) {
	m.register <- client
}

func (m *Manager) Unregister(client types.WebSocketClientInterface) {
	m.unregister <- client
}

func (m *Manager) SendToGame(gameID string, message []byte) {
	m.mu.RLock()
	connections := m.gameConnections[gameID]
	m.mu.RUnlock()

	for conn := range connections {
		if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
			log.Printf("Error sending message: %v", err)
		}
	}
}

func (m *Manager) Run() {
	for {
		select {
		case <-m.shutdown:
			return
		case client := <-m.register:
			m.handleRegister(client)
		case client := <-m.unregister:
			m.handleUnregister(client)
		}
	}
}

func (m *Manager) Stop() {
	close(m.shutdown)
}

func (m *Manager) handleRegister(client types.WebSocketClientInterface) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.gameConnections[client.GetGameID()]; !exists {
		m.gameConnections[client.GetGameID()] = make(map[*websocket.Conn]bool)
	}
	if wsClient, ok := client.(*Client); ok {
		m.gameConnections[client.GetGameID()][wsClient.Socket] = true
	}
}

func (m *Manager) handleUnregister(client types.WebSocketClientInterface) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if wsClient, ok := client.(*Client); ok {
		if conns, exists := m.gameConnections[client.GetGameID()]; exists {
			delete(conns, wsClient.Socket)
		}
	}
}

func (m *Manager) HandleConnection(w http.ResponseWriter, r *http.Request, gameID string, playerID string) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}

	client := NewClient(playerID, gameID, conn)
	m.Register(client)

	// Handle connection in goroutines
	go client.Read()
	go client.Write()
}
