package services

import (
	"taboo-game/models"
	"taboo-game/services"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMatchService(t *testing.T) {
	t.Run("GetMatch", func(t *testing.T) {
		gameSvc := services.NewGameService()
		svc := services.NewMatchService(gameSvc)
		gameID := "game-1"
		matchID := "match-1"
		teamAssignments := map[string][]string{
			"teamA": {"player1", "player2"},
			"teamB": {"player3", "player4"},
		}

		// Create a match first
		match, err := svc.StartMatch(gameID, matchID, teamAssignments)
		assert.NoError(t, err)
		assert.NotNil(t, match)

		// Get the match
		retrievedMatch, err := svc.GetMatch(gameID, matchID)
		assert.NoError(t, err)
		assert.Equal(t, matchID, retrievedMatch.ID)
		assert.Equal(t, gameID, retrievedMatch.GameID)
	})

	t.Run("StartMatch", func(t *testing.T) {
		gameSvc := services.NewGameService()
		svc := services.NewMatchService(gameSvc)
		gameID := "game-1"
		matchID := "match-1"
		teamAssignments := map[string][]string{
			"teamA": {"player1", "player2"},
			"teamB": {"player3", "player4"},
		}

		match, err := svc.StartMatch(gameID, matchID, teamAssignments)

		assert.NoError(t, err)
		assert.Equal(t, matchID, match.ID)
		assert.Equal(t, gameID, match.GameID)
		assert.Equal(t, models.MatchStatusPending, match.Status)
		assert.True(t, match.TeamATurn)
	})

	t.Run("CreateStage", func(t *testing.T) {
		gameSvc := services.NewGameService()
		svc := services.NewMatchService(gameSvc)
		gameID := "game-1"
		matchID := "match-1"
		teamAssignments := map[string][]string{}

		// Create a match first
		_, _ = svc.StartMatch(gameID, matchID, teamAssignments)

		stageDetails := models.MatchStageDetails{
			ActiveTeamID:   "team-1",
			SpottingTeamID: "team-2",
			ClueGivers:     []string{"player1"},
			Guessers:       []string{"player2"},
			Spotters:       []string{"player3"},
		}

		stage, err := svc.CreateStage(gameID, matchID, stageDetails)

		assert.NoError(t, err)
		assert.NotEmpty(t, stage.ID)
		assert.Equal(t, matchID, stage.MatchID)
		assert.Equal(t, "team-1", stage.ActiveTeamID)
		assert.Equal(t, "pending", stage.Status)
	})
}
