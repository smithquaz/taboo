package services_test

import (
	"taboo-game/models"
	"taboo-game/services"
	"taboo-game/tests/mocks"
	"taboo-game/websocket"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setupMatchService(t *testing.T) (*services.MatchService, *models.MatchDetails) {
	wsManager := websocket.NewManager()
	mockGameService := &mocks.MockGameService{
		GetGameFunc: func(gameID string) (*models.Game, error) {
			return &models.Game{
				ID:     "test-game",
				Status: "in_progress",
			}, nil
		},
	}
	ms := services.NewMatchService(mockGameService, wsManager)

	// Create and store a test match
	match := createTestMatch(t)
	ms.StoreMatch(match)

	return ms, match
}

func createTestMatch(t *testing.T) *models.MatchDetails {
	return &models.MatchDetails{
		ID:           "test-match",
		GameID:       "test-game",
		Status:       models.MatchStatusPending,
		TeamAPlayers: []string{"p1", "p2", "p3"},
		TeamBPlayers: []string{"p4", "p5", "p6", "p7"},
		TeamAScore:   0,
		TeamBScore:   0,
		CurrentStage: &models.MatchStage{
			ID:         "test-stage",
			MatchID:    "test-match",
			Status:     "active",
			TeamAScore: 0,
			TeamBScore: 0,
		},
	}
}

func containsPlayer(players []string, playerID string) bool {
	for _, p := range players {
		if p == playerID {
			return true
		}
	}
	return false
}

func TestGetMatch(t *testing.T) {
	ms, _ := setupMatchService(t)
	gameID := "game-1"
	matchID := "match-1"

	// Test getting non-existent match
	match, err := ms.GetMatch(gameID, matchID)
	assert.Error(t, err)
	assert.Nil(t, match)

	// Create a match first
	teamAssignments := map[string][]string{
		"teamA": {"player1", "player2", "player3"},
		"teamB": {"player4", "player5", "player6", "player7"},
	}
	match, err = ms.StartMatch(gameID, matchID, teamAssignments)
	assert.NoError(t, err)
	assert.NotNil(t, match)

	// Test getting existing match
	retrievedMatch, err := ms.GetMatch(gameID, matchID)
	assert.NoError(t, err)
	assert.NotNil(t, retrievedMatch)
	assert.Equal(t, matchID, retrievedMatch.ID)
}

func TestStartMatch(t *testing.T) {
	ms, _ := setupMatchService(t)
	gameID := "game-1"
	matchID := "match-1"

	t.Run("valid team assignments", func(t *testing.T) {
		teamAssignments := map[string][]string{
			"teamA": {"player1", "player2", "player3"},
			"teamB": {"player4", "player5", "player6", "player7"},
		}
		match, err := ms.StartMatch(gameID, matchID, teamAssignments)
		assert.NoError(t, err)
		assert.NotNil(t, match)
		assert.Equal(t, models.MatchStatusPending, match.Status)
		assert.True(t, match.TeamATurn)
	})

	t.Run("invalid team size", func(t *testing.T) {
		teamAssignments := map[string][]string{
			"teamA": {"player1"},
			"teamB": {"player2", "player3"},
		}
		match, err := ms.StartMatch(gameID, matchID, teamAssignments)
		assert.Error(t, err)
		assert.Nil(t, match)
	})

	t.Run("unbalanced teams", func(t *testing.T) {
		teamAssignments := map[string][]string{
			"teamA": {"player1", "player2"},
			"teamB": {"player3", "player4", "player5", "player6"},
		}
		match, err := ms.StartMatch(gameID, matchID, teamAssignments)
		assert.Error(t, err)
		assert.Nil(t, match)
	})
}

func TestProcessGuessAttempt(t *testing.T) {
	ms, match := setupMatchService(t)

	attempt := &models.GuessAttempt{
		CardID:  "test-card",
		Correct: true,
		TeamID:  "teamA",
		StageID: match.CurrentStage.ID,
	}

	err := ms.ProcessGuessAttempt(attempt)
	assert.NoError(t, err)
	assert.Equal(t, models.PointsCorrectGuess, match.CurrentStage.TeamAScore)
	assert.Equal(t, models.PointsCorrectGuess, match.TeamAScore)

	// Test violation scoring
	attempt.Correct = false
	attempt.Violation = true
	attempt.TeamID = "teamA"
	err = ms.ProcessGuessAttempt(attempt)
	assert.NoError(t, err)

	assert.Equal(t, models.PointsViolationCatch, match.CurrentStage.TeamBScore)
	assert.Equal(t, models.PointsViolationCatch, match.TeamBScore)
}

func TestFinalizeStageScores(t *testing.T) {
	ms, match := setupMatchService(t)

	err := ms.FinalizeStageScores(match.CurrentStage.ID)
	assert.NoError(t, err)

	// Team A has 3 players, should get base points
	assert.Equal(t, models.BasePointsTeamOfThree, match.TeamAScore)
	assert.Equal(t, models.BasePointsTeamOfThree, match.CurrentStage.TeamAScore)
}

func TestSwitchTeam(t *testing.T) {
	ms, _ := setupMatchService(t)
	gameID := "game-1"
	matchID := "match-1"

	// Start with balanced teams (3v4)
	teamAssignments := map[string][]string{
		"teamA": {"player1", "player2", "player3"},
		"teamB": {"player4", "player5", "player6", "player7"},
	}

	match, err := ms.StartMatch(gameID, matchID, teamAssignments)
	assert.NoError(t, err)

	t.Run("valid switch - balancing teams", func(t *testing.T) {
		// Switch from larger team (B) to smaller team (A)
		match, err = ms.SwitchTeam(matchID, "player7")
		assert.NoError(t, err)
		assert.Equal(t, 4, len(match.TeamAPlayers))
		assert.Equal(t, 3, len(match.TeamBPlayers))
		assert.True(t, containsPlayer(match.TeamAPlayers, "player7"))
	})

	t.Run("invalid switch - would create imbalance", func(t *testing.T) {
		// Try to switch from smaller team (B) to larger team (A)
		// This would make teams 5v2, which is invalid
		match, err = ms.SwitchTeam(matchID, "player4")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "switch would create team imbalance")
	})

	t.Run("invalid switch - player not found", func(t *testing.T) {
		match, err = ms.SwitchTeam(matchID, "nonexistent")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "player not found in any team")
	})
}
