package websocket_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"taboo-game/websocket"
	"testing"

	gorilla "github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

func TestWebSocketManager(t *testing.T) {
	t.Run("Client Connection and Message Broadcasting", func(t *testing.T) {
		manager := websocket.NewManager()
		go manager.Start()

		// Create test server
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			upgrader := gorilla.Upgrader{}
			conn, err := upgrader.Upgrade(w, r, nil)
			assert.NoError(t, err)

			client := &websocket.Client{
				ID:      "test-client",
				GameID:  "test-game",
				Socket:  conn,
				Send:    make(chan []byte, 256),
				Manager: manager,
			}

			manager.Register(client)
			go client.Read()
			go client.Write()
		}))
		defer server.Close()

		// Connect test client
		wsURL := "ws" + strings.TrimPrefix(server.URL, "http")
		conn, _, err := gorilla.DefaultDialer.Dial(wsURL, nil)
		assert.NoError(t, err)
		defer conn.Close()

		// Test sending message
		testMessage := []byte("test message")
		manager.SendToGame("test-game", testMessage)

		// Wait for message to be received
		_, message, err := conn.ReadMessage()
		assert.NoError(t, err)
		assert.Equal(t, testMessage, message)
	})
}
