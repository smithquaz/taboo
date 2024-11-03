package websocket

import (
	"bytes"
	"context"
	"log"
	"taboo-game/websocket"
	"testing"
	"time"

	"taboo-game/types"
	"github.com/stretchr/testify/assert"
)

// Add debug logger
var debugLog = log.New(log.Writer(), "[DEBUG] ", log.Ltime|log.Lshortfile)

// MockClient implements WebSocketClientInterface
type MockClient struct {
	ID     string
	GameID string
	Send   chan []byte
}

func NewMockClient(id, gameID string) *MockClient {
	return &MockClient{
		ID:     id,
		GameID: gameID,
		Send:   make(chan []byte, 256),
	}
}

func (m *MockClient) GetID() string {
	return m.ID
}

func (m *MockClient) GetGameID() string {
	return m.GameID
}

func (m *MockClient) Read()  {}
func (m *MockClient) Write() {}

// MockGameEventsHandler implements websocket.GameEventsHandler
type MockGameEventsHandler struct {
	StartStageFunc      func(gameID string, stageNum int) error
	HandleClueFunc      func(gameID string, playerID string, clue string) error
	HandleGuessFunc     func(gameID string, playerID string, guess string) error
	HandleViolationFunc func(gameID string, reporterID string, violationType string) error
}

// Implement all interface methods
func (h *MockGameEventsHandler) StartStage(gameID string, stageNum int) error {
	if h.StartStageFunc != nil {
		return h.StartStageFunc(gameID, stageNum)
	}
	return nil
}

func (h *MockGameEventsHandler) HandleClue(gameID string, playerID string, clue string) error {
	if h.HandleClueFunc != nil {
		return h.HandleClueFunc(gameID, playerID, clue)
	}
	return nil
}

func (h *MockGameEventsHandler) HandleGuess(gameID string, playerID string, guess string) error {
	if h.HandleGuessFunc != nil {
		return h.HandleGuessFunc(gameID, playerID, guess)
	}
	return nil
}

func (h *MockGameEventsHandler) HandleViolation(gameID string, reporterID string, violationType string) error {
	if h.HandleViolationFunc != nil {
		return h.HandleViolationFunc(gameID, reporterID, violationType)
	}
	return nil
}

// Update the test to use the new mock
func TestWebSocketManager_Unit(t *testing.T) {
	t.Run("Register clients", func(t *testing.T) {
		mockEvents := &MockGameEventsHandler{}
		manager := websocket.NewManager(mockEvents)
		managerDone := make(chan bool)

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		go func() {
			manager.Run()
			managerDone <- true
		}()

		client := NewMockClient("1", "game1")
		manager.Register(client)

		// Give time for registration
		time.Sleep(100 * time.Millisecond)

		messageReceived := make(chan bool)
		go func() {
			select {
			case msg := <-client.Send:
				assert.Equal(t, []byte("test"), msg)
				messageReceived <- true
			case <-time.After(500 * time.Millisecond):
				debugLog.Printf("Timeout waiting for message")
				messageReceived <- false
			}
		}()

		manager.SendToGame("game1", []byte("test"))

		select {
		case success := <-messageReceived:
			assert.True(t, success, "Message should have been received")
		case <-ctx.Done():
			t.Fatal("Test timed out waiting for message")
		}

		// Graceful shutdown
		manager.Stop()
		select {
		case <-managerDone:
			debugLog.Printf("Manager stopped gracefully")
		case <-time.After(time.Second):
			t.Fatal("Manager failed to shut down")
		}
	})

	t.Run("SendToGame sends to correct clients", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		mockEvents := &MockGameEventsHandler{}
		manager := websocket.NewManager(mockEvents)
		managerDone := make(chan bool)
		go func() {
			manager.Run()
			managerDone <- true
		}()

		// Create and register clients as WebSocketClientInterface
		clients := []types.WebSocketClientInterface{
			NewMockClient("1", "game1"),
			NewMockClient("2", "game1"),
			NewMockClient("3", "game2"),
		}

		for _, client := range clients {
			manager.Register(client)
			debugLog.Printf("Registered client %s for game %s", client.GetID(), client.GetGameID())
		}

		time.Sleep(100 * time.Millisecond)

		results := make(chan struct {
			clientID string
			received bool
		}, len(clients))

		// Monitor each client's Send channel
		for _, client := range clients {
			go func(c types.WebSocketClientInterface) {
				if mockClient, ok := c.(*MockClient); ok {
					select {
					case msg := <-mockClient.Send:
						debugLog.Printf("Client %s received message", c.GetID())
						results <- struct {
							clientID string
							received bool
						}{c.GetID(), bytes.Equal(msg, []byte("test message"))}
					case <-time.After(500 * time.Millisecond):
						debugLog.Printf("Client %s timed out waiting for message", c.GetID())
						results <- struct {
							clientID string
							received bool
						}{c.GetID(), false}
					}
				}
			}(client)
		}

		debugLog.Printf("Sending test message to game1")
		manager.SendToGame("game1", []byte("test message"))

		receivedMap := make(map[string]bool)
		for i := 0; i < len(clients); i++ {
			select {
			case result := <-results:
				receivedMap[result.clientID] = result.received
			case <-ctx.Done():
				t.Fatal("Test timed out waiting for messages")
			}
		}

		// Verify results
		assert.True(t, receivedMap["1"], "Client1 should receive the message")
		assert.True(t, receivedMap["2"], "Client2 should receive the message")
		assert.False(t, receivedMap["3"], "Client3 should not receive the message")

		// Graceful shutdown
		debugLog.Printf("Starting graceful shutdown")
		manager.Stop()
		select {
		case <-managerDone:
			debugLog.Printf("Manager stopped gracefully")
		case <-time.After(time.Second):
			t.Fatal("Manager failed to shut down")
		}
	})
}

func TestReconnection(t *testing.T) {
	t.Run("Client reconnection with message replay", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		mockEvents := &MockGameEventsHandler{}
		manager := websocket.NewManager(mockEvents)
		go manager.Run()
		defer manager.Stop()

		// Create and register initial client
		client1 := NewMockClient("1", "game1")
		manager.Register(client1)
		debugLog.Printf("Registered initial client")

		// Give time for registration to complete
		time.Sleep(100 * time.Millisecond)

		// Send first message
		message1 := []byte("message before disconnect")
		manager.SendToGame("game1", message1)
		debugLog.Printf("Sent pre-disconnect message")

		// Verify first message received
		select {
		case msg := <-client1.Send:
			assert.Equal(t, message1, msg)
		case <-time.After(100 * time.Millisecond):
			t.Fatal("Failed to receive first message")
		}

		// Simulate disconnect
		manager.Unregister(client1)
		debugLog.Printf("Client disconnected")

		// Send message during disconnect
		message2 := []byte("message during disconnect")
		manager.SendToGame("game1", message2)
		debugLog.Printf("Sent message during disconnect")

		// Give time for message to be queued
		time.Sleep(100 * time.Millisecond)

		// Reconnect with same ID
		client1Reconnect := NewMockClient("1", "game1")
		manager.Register(client1Reconnect)
		debugLog.Printf("Client reconnected")

		// Wait for replayed messages
		receivedMessages := make([][]byte, 0)
		messageTimer := time.NewTimer(500 * time.Millisecond)

	CollectMessages:
		for {
			select {
			case msg := <-client1Reconnect.Send:
				debugLog.Printf("Received replayed message: %s", string(msg))
				receivedMessages = append(receivedMessages, msg)
				if len(receivedMessages) == 2 {
					break CollectMessages
				}
			case <-messageTimer.C:
				debugLog.Printf("Timeout waiting for messages. Received %d messages", len(receivedMessages))
				break CollectMessages
			case <-ctx.Done():
				t.Fatal("Test timed out")
			}
		}

		// Verify messages
		assert.Equal(t, 2, len(receivedMessages), "Should receive both messages")
		if len(receivedMessages) >= 2 {
			assert.Equal(t, message1, receivedMessages[0], "First message should match")
			assert.Equal(t, message2, receivedMessages[1], "Second message should match")
		}
	})
}
