package routes

import (
	"github.com/gin-gonic/gin"
	"taboo-game/handlers"
)

type MatchRoutes struct {
	matchHandler *handlers.MatchHandler
}

func NewMatchRoutes(matchHandler *handlers.MatchHandler) *MatchRoutes {
	return &MatchRoutes{
		matchHandler: matchHandler,
	}
}

func (r *MatchRoutes) RegisterRoutes(rg *gin.Engine) {
	matches := rg.Group("/api/v1/games/:gameId/matches")
	{
		matches.GET("/:matchId", r.matchHandler.GetMatch)
		matches.PUT("/:matchId/start", r.matchHandler.StartMatch)
		matches.PUT("/:matchId/end", r.matchHandler.EndMatch)
		matches.POST("/:matchId/stages", r.matchHandler.CreateStage)
	}
}
