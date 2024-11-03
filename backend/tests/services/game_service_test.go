package services

import (
	"github.com/stretchr/testify/assert"
	"taboo-game/models"
	"taboo-game/services"
	"testing"
)

func TestGameService(t *testing.T) {
	t.Run("CreateGame", func(t *testing.T) {
		svc := services.NewGameService()

		game, err := svc.CreateGame(4)

		assert.NoError(t, err)
		assert.NotEmpty(t, game.ID)
		assert.Equal(t, models.GameStatus("waiting"), game.Status)
		assert.Len(t, game.Teams, 2)
		assert.Equal(t, 4, game.Teams[0].Size)
		assert.Equal(t, 4, game.Teams[1].Size)
	})

	t.Run("AddPlayer", func(t *testing.T) {
		svc := services.NewGameService()
		game, _ := svc.CreateGame(4)

		player, err := svc.AddPlayer(game.ID, "TestPlayer")

		assert.NoError(t, err)
		assert.NotEmpty(t, player.ID)
		assert.Equal(t, "TestPlayer", player.Name)

		// Verify player was added to a team
		updatedGame, _ := svc.GetGame(game.ID)
		assert.Len(t, updatedGame.Teams[0].Players, 1)
	})

	t.Run("StartGame", func(t *testing.T) {
		svc := services.NewGameService()
		game, _ := svc.CreateGame(2)

		// Add players to fill teams
		svc.AddPlayer(game.ID, "Player1")
		svc.AddPlayer(game.ID, "Player2")
		svc.AddPlayer(game.ID, "Player3")
		svc.AddPlayer(game.ID, "Player4")

		startedGame, err := svc.StartGame(game.ID)

		assert.NoError(t, err)
		assert.Equal(t, models.GameStatusInProgress, startedGame.Status)
		assert.Len(t, startedGame.Matches, 3)
	})

	t.Run("EndGame", func(t *testing.T) {
		svc := services.NewGameService()
		game, _ := svc.CreateGame(2)

		// Setup game with players and start it
		svc.AddPlayer(game.ID, "Player1")
		svc.AddPlayer(game.ID, "Player2")
		svc.AddPlayer(game.ID, "Player3")
		svc.AddPlayer(game.ID, "Player4")
		svc.StartGame(game.ID)

		endedGame, err := svc.EndGame(game.ID)

		assert.NoError(t, err)
		assert.Equal(t, models.GameStatusCompleted, endedGame.Status)
	})

	t.Run("GetGame_NotFound", func(t *testing.T) {
		svc := services.NewGameService()

		game, err := svc.GetGame("non-existent-id")

		assert.Error(t, err)
		assert.Nil(t, game)
		assert.Equal(t, "game not found", err.Error())
	})
}
