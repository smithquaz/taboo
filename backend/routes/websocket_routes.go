package routes

import (
	"github.com/gin-gonic/gin"
	"taboo-game/websocket"
)

func SetupWebSocketRoutes(router *gin.Engine, wsManager *websocket.Manager) {
	router.GET("/ws/:gameId/:playerId", func(c *gin.Context) {
		gameID := c.Param("gameId")
		playerID := c.Param("playerId")
		wsManager.HandleConnection(c.Writer, c.Request, gameID, playerID)
	})
}
