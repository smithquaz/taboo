package websocket

import (
	"sync"
)

type mockGameEvents struct {
	eventReceived chan bool
	mu            sync.Mutex
	calls         []struct {
		gameID   string
		playerID string
		event    interface{}
	}
}

func NewMockGameEvents() *mockGameEvents {
	return &mockGameEvents{
		eventReceived: make(chan bool, 1),
		calls: make([]struct {
			gameID   string
			playerID string
			event    interface{}
		}, 0),
	}
}

func (m *mockGameEvents) HandleGameEvent(gameID, playerID string, event interface{}) error {
	m.mu.Lock()
	m.calls = append(m.calls, struct {
		gameID   string
		playerID string
		event    interface{}
	}{gameID, playerID, event})
	m.mu.Unlock()

	select {
	case m.eventReceived <- true:
	default:
	}
	return nil
}
