package websocket

import (
	"github.com/gorilla/websocket"
)

type Client struct {
	ID     string
	GameID string
	Socket *websocket.Conn
	Send   chan []byte
}

func NewClient(id, gameID string, socket *websocket.Conn) *Client {
	return &Client{
		ID:     id,
		GameID: gameID,
		Socket: socket,
		Send:   make(chan []byte, 256),
	}
}

func (c *Client) GetID() string {
	return c.ID
}

func (c *Client) GetGameID() string {
	return c.GameID
}

func (c *Client) GetSocket() *websocket.Conn {
	return c.Socket
}

func (c *Client) Read() {
	defer c.Socket.Close()

	for {
		_, message, err := c.Socket.ReadMessage()
		if err != nil {
			break
		}

		c.Send <- message
	}
}

func (c *Client) Write() {
	for message := range c.Send {
		c.Socket.WriteMessage(websocket.TextMessage, message)
	}
}
