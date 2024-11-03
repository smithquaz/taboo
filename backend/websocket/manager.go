package websocket

import (
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

type Client struct {
	ID      string
	GameID  string
	Socket  *websocket.Conn
	Send    chan []byte
	Manager *Manager
}

type Manager struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	sync.RWMutex
	// Maps gameID to a list of clients
	gameClients map[string]map[*Client]bool
}

func NewManager() *Manager {
	return &Manager{
		clients:     make(map[*Client]bool),
		broadcast:   make(chan []byte),
		register:    make(chan *Client),
		unregister:  make(chan *Client),
		gameClients: make(map[string]map[*Client]bool),
	}
}

// Register adds a new client to the manager
func (m *Manager) Register(client *Client) {
	m.register <- client
}

func (m *Manager) Start() {
	for {
		select {
		case client := <-m.register:
			m.Lock()
			m.clients[client] = true
			if _, ok := m.gameClients[client.GameID]; !ok {
				m.gameClients[client.GameID] = make(map[*Client]bool)
			}
			m.gameClients[client.GameID][client] = true
			m.Unlock()

		case client := <-m.unregister:
			m.Lock()
			if _, ok := m.clients[client]; ok {
				delete(m.clients, client)
				delete(m.gameClients[client.GameID], client)
				close(client.Send)
			}
			m.Unlock()

		case message := <-m.broadcast:
			m.Lock()
			for client := range m.clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(m.clients, client)
					delete(m.gameClients[client.GameID], client)
				}
			}
			m.Unlock()
		}
	}
}

func (m *Manager) SendToGame(gameID string, message []byte) {
	m.Lock()
	defer m.Unlock()

	if clients, ok := m.gameClients[gameID]; ok {
		for client := range clients {
			select {
			case client.Send <- message:
			default:
				close(client.Send)
				delete(m.clients, client)
				delete(m.gameClients[client.GameID], client)
			}
		}
	}
}

func (c *Client) Read() {
	defer func() {
		c.Manager.unregister <- c
		c.Socket.Close()
	}()

	for {
		_, message, err := c.Socket.ReadMessage()
		if err != nil {
			log.Printf("Error reading from websocket: %v", err)
			break
		}
		// Handle incoming messages
		log.Printf("Received message from client %s: %s", c.ID, string(message))
	}
}

func (c *Client) Write() {
	defer func() {
		c.Socket.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				c.Socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			err := c.Socket.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				log.Printf("Error writing to websocket: %v", err)
				return
			}
		}
	}
}
