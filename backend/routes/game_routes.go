package routes

import (
	"taboo-game/handlers"

	"github.com/gin-gonic/gin"
)

type GameRoutes struct {
	gameHandler *handlers.GameHandler
}

func NewGameRoutes(gameHandler *handlers.GameHandler) *GameRoutes {
	return &GameRoutes{
		gameHandler: gameHandler,
	}
}

func (r *GameRoutes) RegisterRoutes(router *gin.Engine) {
	api := router.Group("/api/v1/games")
	{
		api.POST("/", r.gameHandler.CreateGame)
		api.POST("/:gameId/join", r.gameHandler.JoinGame)
		api.GET("/:gameId", r.gameHandler.GetGame)
		api.PUT("/:gameId/start", r.gameHandler.StartGame)
		api.PUT("/:gameId/end", r.gameHandler.EndGame)
	}
}
