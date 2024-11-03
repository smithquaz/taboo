package mocks

import (
	"net/http"
	"taboo-game/types"
)

type MockWebSocketManager struct {
	SendToGameFunc func(gameID string, message []byte)
	RegisterFunc   func(client types.WebSocketClientInterface)
	UnregisterFunc func(client types.WebSocketClientInterface)
}

func (m *MockWebSocketManager) SendToGame(gameID string, message []byte) {
	if m.SendToGameFunc != nil {
		m.SendToGameFunc(gameID, message)
	}
}

func (m *MockWebSocketManager) Register(client types.WebSocketClientInterface) {
	if m.RegisterFunc != nil {
		m.RegisterFunc(client)
	}
}

func (m *MockWebSocketManager) Unregister(client types.WebSocketClientInterface) {
	if m.UnregisterFunc != nil {
		m.UnregisterFunc(client)
	}
}

func (m *MockWebSocketManager) HandleConnection(w http.ResponseWriter, r *http.Request, gameID, playerID string) {
}
func (m *MockWebSocketManager) Run() {}
