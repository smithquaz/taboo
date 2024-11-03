package handlers

import (
	"github.com/gin-gonic/gin"
)

type PlayerHandler struct {
	// Dependencies would go here
	// playerService *services.PlayerService
	// db           *database.DB
}

func NewPlayerHandler() *PlayerHandler {
	return &PlayerHandler{}
}

func (h *PlayerHandler) CreatePlayer(c *gin.Context) {
	// TODO: Parse request body into player model
	// TODO: Call service to create player
	c.JSON(200, gin.H{"message": "Player created"})
}

func (h *PlayerHandler) GetPlayer(c *gin.Context) {
	id := c.Param("id")
	// TODO: Call service to get player by id
	c.JSON(200, gin.H{"message": "Get player " + id})
}

func (h *PlayerHandler) UpdatePlayer(c *gin.Context) {
	id := c.Param("id")
	// TODO: Parse request body and update player
	c.JSON(200, gin.H{"message": "Update player " + id})
}

func (h *PlayerHandler) DeletePlayer(c *gin.Context) {
	id := c.Param("id")
	// TODO: Call service to delete player
	c.JSON(200, gin.H{"message": "Delete player " + id})
}
