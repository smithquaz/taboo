package handlers

import (
	"net/http"
	"taboo-game/models"
	"taboo-game/services"

	"github.com/gin-gonic/gin"
)

type MatchHandler struct {
	matchService services.MatchServiceInterface
}

func NewMatchHandler(matchService services.MatchServiceInterface) *MatchHandler {
	return &MatchHandler{
		matchService: matchService,
	}
}

func (h *MatchHandler) StartMatch(c *gin.Context) {
	gameID := c.Param("gameId")
	matchID := c.Param("matchId")

	var req struct {
		TeamAssignments map[string][]string `json:"teamAssignments"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	match, err := h.matchService.StartMatch(gameID, matchID, req.TeamAssignments)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, match)
}

func (h *MatchHandler) ScorePoint(c *gin.Context) {
	matchID := c.Param("matchId")
	var req struct {
		IsTeamA bool `json:"isTeamA"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	match, err := h.matchService.ScorePoint(matchID, req.IsTeamA)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, match)
}

func (h *MatchHandler) EndMatch(c *gin.Context) {
	gameID := c.Param("gameId")
	matchID := c.Param("matchId")

	match, err := h.matchService.EndMatch(gameID, matchID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, match)
}

func (h *MatchHandler) CreateStage(c *gin.Context) {
	gameID := c.Param("gameId")
	matchID := c.Param("matchId")

	var req models.MatchStageDetails
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	stage, err := h.matchService.CreateStage(gameID, matchID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, stage)
}

func (h *MatchHandler) GetMatch(c *gin.Context) {
	gameID := c.Param("gameId")
	matchID := c.Param("matchId")

	match, err := h.matchService.GetMatch(gameID, matchID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, match)
}
