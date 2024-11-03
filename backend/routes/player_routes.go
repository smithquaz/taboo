package routes

import (
	"taboo-game/handlers"

	"github.com/gin-gonic/gin"
)

type PlayerRoutes struct {
	playerHandler *handlers.PlayerHandler
}

func NewPlayerRoutes(playerHandler *handlers.PlayerHandler) *PlayerRoutes {
	return &PlayerRoutes{
		playerHandler: playerHandler,
	}
}

func (r *PlayerRoutes) RegisterRoutes(rg *gin.Engine) {
	players := rg.Group("/api/v1/players")
	{
		players.POST("/", r.playerHandler.CreatePlayer)
		players.GET("/:id", r.playerHandler.GetPlayer)
		players.PUT("/:id", r.playerHandler.UpdatePlayer)
		players.DELETE("/:id", r.playerHandler.DeletePlayer)
	}
}
