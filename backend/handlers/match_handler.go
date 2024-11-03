package handlers

import (
	"net/http"
	"taboo-game/services"

	"github.com/gin-gonic/gin"
)

type MatchHandler struct {
	matchService *services.MatchService
}

func NewMatchHandler(matchService *services.MatchService) *MatchHandler {
	return &MatchHandler{
		matchService: matchService,
	}
}

func (h *MatchHandler) GetMatch(c *gin.Context) {
	gameID := c.Param("gameId")
	matchID := c.Param("matchId")

	match, err := h.matchService.GetMatch(gameID, matchID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, match)
}

func (h *MatchHandler) StartMatch(c *gin.Context) {
	gameID := c.Param("gameId")
	matchID := c.Param("matchId")

	var req struct {
		TeamAssignments map[string][]string `json:"teamAssignments"` // teamID -> playerIDs
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

// TODO: Implement this
func (h *MatchHandler) CreateStage(c *gin.Context) {
	// gameID := c.Param("gameId")
	// matchID := c.Param("matchId")

	var req struct {
		ActiveTeamID   string   `json:"activeTeamId"`
		SpottingTeamID string   `json:"spottingTeamId"`
		ClueGivers     []string `json:"clueGivers"`
		Guessers       []string `json:"guessers"`
		Spotters       []string `json:"spotters"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// stage, err := h.matchService.CreateStage(gameID, matchID, req)
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 	return
	// }
	// c.JSON(http.StatusCreated, stage)
}
