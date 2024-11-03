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

func (r *GameRoutes) RegisterRoutes(rg *gin.Engine) {
	games := rg.Group("/api/v1/games")
	{
		games.POST("/", r.gameHandler.CreateGame)
		games.POST("/:id/join", r.gameHandler.JoinGame)
		games.GET("/:id", r.gameHandler.GetGame)
		games.PUT("/:id/start", r.gameHandler.StartGame)
		games.PUT("/:id/end", r.gameHandler.EndGame)
	}
}
