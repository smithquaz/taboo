package handlers

import (
	"github.com/gin-gonic/gin"
	"taboo-game/types"
)

type WebSocketHandler struct {
	manager    types.WebSocketManagerInterface
	gameEvents types.GameEventsServiceInterface
}

func NewWebSocketHandler(manager types.WebSocketManagerInterface, gameEvents types.GameEventsServiceInterface) *WebSocketHandler {
	return &WebSocketHandler{
		manager:    manager,
		gameEvents: gameEvents,
	}
}

func (h *WebSocketHandler) HandleConnection(c *gin.Context) {
	gameID := c.Param("gameId")
	playerID := c.Param("playerId")
	h.manager.HandleConnection(c.Writer, c.Request, gameID, playerID)
}
