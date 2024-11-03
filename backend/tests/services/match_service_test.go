package services

import (
	"taboo-game/models"
	"taboo-game/services"
	"taboo-game/tests/mocks"
	"taboo-game/websocket"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Helper function for the test
func containsPlayer(players []string, playerID string) bool {
	for _, p := range players {
		if p == playerID {
			return true
		}
	}
	return false
}

func TestMatchService(t *testing.T) {
	t.Run("GetMatch", func(t *testing.T) {
		mockGameSvc := &mocks.MockGameService{
			GetGameFunc: func(gameID string) (*models.Game, error) {
				return &models.Game{
					ID:        gameID,
					Status:    "in_progress",
					CreatedAt: time.Now(),
					Teams: []models.Team{
						{
							ID:   "teamA",
							Size: 3,
						},
						{
							ID:   "teamB",
							Size: 4,
						},
					},
				}, nil
			},
		}

		wsManager := websocket.NewManager()
		svc := services.NewMatchService(mockGameSvc, wsManager)
		gameID := "game-1"
		matchID := "match-1"
		teamAssignments := map[string][]string{
			"teamA": {"player1", "player2", "player3"},
			"teamB": {"player4", "player5", "player6", "player7"},
		}

		// Create a match
		match, err := svc.StartMatch(gameID, matchID, teamAssignments)
		assert.NoError(t, err)
		assert.NotNil(t, match)

		// Get the match
		retrievedMatch, err := svc.GetMatch(gameID, matchID)
		assert.NoError(t, err)
		assert.NotNil(t, retrievedMatch)
		assert.Equal(t, matchID, retrievedMatch.ID)
		assert.Equal(t, gameID, retrievedMatch.GameID)
		assert.Equal(t, 3, len(retrievedMatch.TeamAPlayers))
		assert.Equal(t, 4, len(retrievedMatch.TeamBPlayers))
	})

	t.Run("StartMatch", func(t *testing.T) {
		mockGameSvc := &mocks.MockGameService{
			GetGameFunc: func(gameID string) (*models.Game, error) {
				return &models.Game{
					ID:        gameID,
					Status:    "in_progress",
					CreatedAt: time.Now(),
					Teams: []models.Team{
						{ID: "teamA", Size: 3},
						{ID: "teamB", Size: 4},
					},
				}, nil
			},
		}

		wsManager := websocket.NewManager()
		svc := services.NewMatchService(mockGameSvc, wsManager)
		gameID := "game-1"
		matchID := "match-1"
		teamAssignments := map[string][]string{
			"teamA": {"player1", "player2", "player3"},
			"teamB": {"player4", "player5", "player6", "player7"},
		}

		match, err := svc.StartMatch(gameID, matchID, teamAssignments)
		assert.NoError(t, err)
		assert.Equal(t, matchID, match.ID)
		assert.Equal(t, gameID, match.GameID)
		assert.Equal(t, models.MatchStatusPending, match.Status)
		assert.True(t, match.TeamATurn)
	})

	t.Run("CreateStage", func(t *testing.T) {
		mockGameSvc := &mocks.MockGameService{
			GetGameFunc: func(gameID string) (*models.Game, error) {
				return &models.Game{
					ID:        gameID,
					Status:    "in_progress",
					CreatedAt: time.Now(),
					Teams: []models.Team{
						{ID: "teamA", Size: 3},
						{ID: "teamB", Size: 4},
					},
				}, nil
			},
		}

		wsManager := websocket.NewManager()
		svc := services.NewMatchService(mockGameSvc, wsManager)
		gameID := "game-1"
		matchID := "match-1"
		teamAssignments := map[string][]string{
			"teamA": {"player1", "player2", "player3"},
			"teamB": {"player4", "player5", "player6", "player7"},
		}

		// Create a match
		_, err := svc.StartMatch(gameID, matchID, teamAssignments)
		assert.NoError(t, err)

		stageDetails := models.MatchStageDetails{
			ActiveTeamID:   "teamA",
			SpottingTeamID: "teamB",
			ClueGivers:     []string{"player1", "player2"},
			Guessers:       []string{"player3"},
			Spotters:       []string{"player4", "player5"},
		}

		stage, err := svc.CreateStage(gameID, matchID, stageDetails)
		assert.NoError(t, err)
		assert.NotEmpty(t, stage.ID)
		assert.Equal(t, matchID, stage.MatchID)
		assert.Equal(t, "teamA", stage.ActiveTeamID)
		assert.Equal(t, "pending", stage.Status)
	})

	t.Run("SwitchTeam", func(t *testing.T) {
		mockGameSvc := &mocks.MockGameService{
			GetGameFunc: func(gameID string) (*models.Game, error) {
				return &models.Game{
					ID:        gameID,
					Status:    "in_progress",
					CreatedAt: time.Now(),
					Teams: []models.Team{
						{ID: "teamA"},
						{ID: "teamB"},
					},
				}, nil
			},
		}

		wsManager := websocket.NewManager()
		svc := services.NewMatchService(mockGameSvc, wsManager)
		gameID := "game-1"
		matchID := "match-1"

		// Start with uneven teams (3v2)
		teamAssignments := map[string][]string{
			"teamA": {"player1", "player2", "player3"},
			"teamB": {"player4", "player5"},
		}

		// Create initial match
		match, err := svc.StartMatch(gameID, matchID, teamAssignments)
		assert.NoError(t, err)
		assert.Equal(t, 3, len(match.TeamAPlayers))
		assert.Equal(t, 2, len(match.TeamBPlayers))

		// Switch a player from team A to B
		match, err = svc.SwitchTeam(matchID, "player3")
		assert.NoError(t, err)
		assert.Equal(t, 2, len(match.TeamAPlayers))
		assert.Equal(t, 3, len(match.TeamBPlayers))
		assert.True(t, containsPlayer(match.TeamBPlayers, "player3"))
	})
}
