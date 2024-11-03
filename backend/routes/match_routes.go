package routes

import (
	"taboo-game/handlers"

	"github.com/gin-gonic/gin"
)

type MatchRoutes struct {
	matchHandler *handlers.MatchHandler
}

func NewMatchRoutes(matchHandler *handlers.MatchHandler) *MatchRoutes {
	return &MatchRoutes{
		matchHandler: matchHandler,
	}
}

func (r *MatchRoutes) RegisterRoutes(router *gin.Engine) {
	api := router.Group("/api/v1")
	{
		matches := api.Group("/matches")
		{
			matches.GET("/:matchId", r.matchHandler.GetMatch)
			matches.POST("/:matchId/start", r.matchHandler.StartMatch)
			matches.POST("/:matchId/guess", r.matchHandler.ProcessGuessAttempt)
			matches.PUT("/:matchId/end", r.matchHandler.EndMatch)
			matches.POST("/:matchId/teams/switch/:playerId", r.matchHandler.SwitchTeam)
		}
	}
}
