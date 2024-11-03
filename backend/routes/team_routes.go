package routes

import (
	"github.com/gin-gonic/gin"
	"taboo-game/handlers"
)

type TeamRoutes struct {
	teamHandler *handlers.TeamHandler
}

func NewTeamRoutes(teamHandler *handlers.TeamHandler) *TeamRoutes {
	return &TeamRoutes{
		teamHandler: teamHandler,
	}
}

func (r *TeamRoutes) RegisterRoutes(rg *gin.Engine) {
	teams := rg.Group("/api/v1/teams")
	{
		teams.POST("/", r.teamHandler.CreateTeam)
		teams.GET("/:id", r.teamHandler.GetTeam)
		teams.PUT("/:id", r.teamHandler.UpdateTeam)
		teams.PUT("/:id/players", r.teamHandler.UpdateTeamPlayers)
		teams.DELETE("/:id", r.teamHandler.DeleteTeam)
	}
}
