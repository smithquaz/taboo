package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"taboo-game/services"
)

type GameHandler struct {
	gameService *services.GameService
}

func NewGameHandler(gameService *services.GameService) *GameHandler {
	return &GameHandler{
		gameService: gameService,
	}
}

func (h *GameHandler) CreateGame(c *gin.Context) {
	var req struct {
		TeamSize int `json:"teamSize" binding:"required,oneof=3 4"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	game, err := h.gameService.CreateGame(req.TeamSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, game)
}

func (h *GameHandler) JoinGame(c *gin.Context) {
	var req struct {
		PlayerName string `json:"playerName" binding:"required"`
	}

	gameID := c.Param("id")
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	player, err := h.gameService.AddPlayer(gameID, req.PlayerName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, player)
}

func (h *GameHandler) GetGame(c *gin.Context) {
	// Implementation
}

func (h *GameHandler) StartGame(c *gin.Context) {
	// Implementation
}

func (h *GameHandler) EndGame(c *gin.Context) {
	// Implementation
}
