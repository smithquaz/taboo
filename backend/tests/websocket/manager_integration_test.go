package websocket

import (
	"context"
	"encoding/json"
	"net/http/httptest"
	"strings"
	"taboo-game/websocket"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	gorilla "github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

// MockGameEvents implements websocket.GameEventsHandler
type MockGameEvents struct {
	StartStageFunc      func(gameID string, stageNum int) error
	HandleClueFunc      func(gameID string, playerID string, clue string) error
	HandleGuessFunc     func(gameID string, playerID string, guess string) error
	HandleViolationFunc func(gameID string, reporterID string, violationType string) error
}

func (m *MockGameEvents) StartStage(gameID string, stageNum int) error {
	if m.StartStageFunc != nil {
		return m.StartStageFunc(gameID, stageNum)
	}
	return nil
}

func (m *MockGameEvents) HandleClue(gameID string, playerID string, clue string) error {
	if m.HandleClueFunc != nil {
		return m.HandleClueFunc(gameID, playerID, clue)
	}
	return nil
}

func (m *MockGameEvents) HandleGuess(gameID string, playerID string, guess string) error {
	if m.HandleGuessFunc != nil {
		return m.HandleGuessFunc(gameID, playerID, guess)
	}
	return nil
}

func (m *MockGameEvents) HandleViolation(gameID string, reporterID string, violationType string) error {
	if m.HandleViolationFunc != nil {
		return m.HandleViolationFunc(gameID, reporterID, violationType)
	}
	return nil
}

func setupTestServer(t *testing.T) (*httptest.Server, *websocket.Manager, *MockGameEvents, chan bool) {
	done := make(chan bool)
	gin.SetMode(gin.TestMode)
	router := gin.New()

	mockEvents := &MockGameEvents{
		StartStageFunc: func(gameID string, stageNum int) error {
			return nil
		},
	}
	manager := websocket.NewManager(mockEvents)

	router.GET("/ws/:gameId/:playerId", func(c *gin.Context) {
		gameID := c.Param("gameId")
		playerID := c.Param("playerId")
		manager.HandleConnection(c.Writer, c.Request, gameID, playerID)
	})

	server := httptest.NewServer(router)
	go func() {
		manager.Run()
		done <- true
	}()

	return server, manager, mockEvents, done
}

func TestWebSocketManager_Integration(t *testing.T) {
	t.Run("Full connection and message cycle", func(t *testing.T) {
		// Create a context with timeout
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		server, manager, _, managerDone := setupTestServer(t)
		defer func() {
			server.Close()
			manager.Stop()
			select {
			case <-managerDone:
				// Clean shutdown
			case <-time.After(time.Second):
				t.Error("Manager failed to shut down")
			}
		}()

		// Connect client with timeout
		wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws/game1/player1"
		dialCtx, dialCancel := context.WithTimeout(ctx, time.Second)
		defer dialCancel()

		conn, _, err := gorilla.DefaultDialer.DialContext(dialCtx, wsURL, nil)
		if err != nil {
			t.Fatalf("Failed to connect: %v", err)
		}
		defer conn.Close()

		// Test sending a stage start message
		message := websocket.Message{
			Type:     websocket.StartStage,
			GameID:   "game1",
			PlayerID: "player1",
			Payload: map[string]interface{}{
				"stage_num": 1,
			},
		}

		// Send message with timeout
		writeCtx, writeCancel := context.WithTimeout(ctx, time.Second)
		defer writeCancel()

		writeComplete := make(chan error)
		go func() {
			data, _ := json.Marshal(message)
			writeComplete <- conn.WriteMessage(gorilla.TextMessage, data)
		}()

		select {
		case err := <-writeComplete:
			assert.NoError(t, err)
		case <-writeCtx.Done():
			t.Fatal("Write timeout")
		}

		// Wait for response with timeout
		readCtx, readCancel := context.WithTimeout(ctx, time.Second)
		defer readCancel()

		readComplete := make(chan struct {
			msg []byte
			err error
		})

		go func() {
			_, msg, err := conn.ReadMessage()
			readComplete <- struct {
				msg []byte
				err error
			}{msg, err}
		}()

		select {
		case result := <-readComplete:
			assert.NoError(t, result.err)
			var response websocket.Message
			err := json.Unmarshal(result.msg, &response)
			assert.NoError(t, err)
			assert.Equal(t, message.GameID, response.GameID)
		case <-readCtx.Done():
			t.Fatal("Read timeout")
		}
	})

	t.Run("Multiple clients in same game", func(t *testing.T) {
		server, manager, _, managerDone := setupTestServer(t)
		defer func() {
			server.Close()
			manager.Stop()
			<-managerDone
		}()

		// Connect two clients to same game
		wsURL1 := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws/game1/player1"
		wsURL2 := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws/game1/player2"

		conn1, _, err1 := gorilla.DefaultDialer.Dial(wsURL1, nil)
		assert.NoError(t, err1)
		defer conn1.Close()

		conn2, _, err2 := gorilla.DefaultDialer.Dial(wsURL2, nil)
		assert.NoError(t, err2)
		defer conn2.Close()

		// Set up timeout context
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		// Send message from one client
		message := websocket.Message{
			Type:     websocket.GiveClue,
			GameID:   "game1",
			PlayerID: "player1",
			Payload: map[string]interface{}{
				"clue": "test clue",
			},
		}

		data, _ := json.Marshal(message)
		err := conn1.WriteMessage(gorilla.TextMessage, data)
		assert.NoError(t, err)

		// Wait for second client to receive message
		messageChan := make(chan []byte)
		errChan := make(chan error)

		go func() {
			_, msg, err := conn2.ReadMessage()
			if err != nil {
				errChan <- err
				return
			}
			messageChan <- msg
		}()

		select {
		case <-ctx.Done():
			t.Fatal("Test timed out waiting for message")
		case err := <-errChan:
			t.Fatalf("Error reading message: %v", err)
		case msg := <-messageChan:
			assert.NotNil(t, msg)
			var response websocket.Message
			err := json.Unmarshal(msg, &response)
			assert.NoError(t, err)
			assert.Equal(t, message.Type, response.Type)
			assert.Equal(t, message.GameID, response.GameID)
		}
	})
}
