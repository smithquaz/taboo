package services

import (
	"errors"
	"taboo-game/models"

	"github.com/google/uuid"
)

type MatchService struct {
	matches     map[string]*models.MatchDetails
	words       []string
	gameService GameServiceInterface
}

func NewMatchService(gameService GameServiceInterface) *MatchService {
	return &MatchService{
		matches:     make(map[string]*models.MatchDetails),
		words:       []string{},
		gameService: gameService,
	}
}

func (s *MatchService) GetMatch(gameID, matchID string) (*models.MatchDetails, error) {
	match, exists := s.matches[matchID]
	if !exists {
		return nil, errors.New("match not found")
	}
	return match, nil
}

func (s *MatchService) StartMatch(gameID, matchID string, teamAssignments map[string][]string) (*models.MatchDetails, error) {
	// First verify the game exists
	game, err := s.gameService.GetGame(gameID)
	if err != nil {
		return nil, errors.New("game not found")
	}

	match, exists := s.matches[matchID]
	if !exists {
		match = &models.MatchDetails{
			ID:     matchID,
			GameID: gameID,
			Status: models.MatchStatusPending,
		}
		s.matches[matchID] = match
	}

	// Validate team assignments
	teamAPlayers, hasTeamA := teamAssignments["teamA"]
	teamBPlayers, hasTeamB := teamAssignments["teamB"]
	if !hasTeamA || !hasTeamB {
		return nil, errors.New("both teams must be assigned")
	}

	// Verify all players exist in the game
	for _, playerID := range teamAPlayers {
		found := false
		for _, team := range game.Teams {
			for _, player := range team.Players {
				if player.ID == playerID {
					found = true
					break
				}
			}
		}
		if !found {
			return nil, errors.New("player not found in game: " + playerID)
		}
	}

	for _, playerID := range teamBPlayers {
		found := false
		for _, team := range game.Teams {
			for _, player := range team.Players {
				if player.ID == playerID {
					found = true
					break
				}
			}
		}
		if !found {
			return nil, errors.New("player not found in game: " + playerID)
		}
	}

	// Update match with team assignments
	match.TeamAPlayers = teamAPlayers
	match.TeamBPlayers = teamBPlayers
	match.Status = models.MatchStatusPending
	match.CurrentWord = s.getNextWord()
	match.TeamATurn = true

	return match, nil
}

func (s *MatchService) ScorePoint(matchID string, isTeamA bool) (*models.MatchDetails, error) {
	match, exists := s.matches[matchID]
	if !exists {
		return nil, errors.New("match not found")
	}

	if match.Status != models.MatchStatusPending {
		return nil, errors.New("match is not in progress")
	}

	if isTeamA {
		match.TeamAScore++
	} else {
		match.TeamBScore++
	}

	return match, nil
}

func (s *MatchService) EndMatch(gameID, matchID string) (*models.MatchDetails, error) {
	match, exists := s.matches[matchID]
	if !exists {
		return nil, errors.New("match not found")
	}

	match.Status = models.MatchStatusCompleted
	return match, nil
}

func (s *MatchService) getNextWord() string {
	return "placeholder"
}

func (s *MatchService) CreateStage(gameID, matchID string, details models.MatchStageDetails) (*models.MatchStage, error) {
	match, exists := s.matches[matchID]
	if !exists {
		return nil, errors.New("match not found")
	}

	if match.Status != models.MatchStatusPending {
		return nil, errors.New("match is not in pending state")
	}

	stage := &models.MatchStage{
		ID:             generateID(),
		MatchID:        matchID,
		ActiveTeamID:   details.ActiveTeamID,
		SpottingTeamID: details.SpottingTeamID,
		ClueGivers:     details.ClueGivers,
		Guessers:       details.Guessers,
		Spotters:       details.Spotters,
		Status:         "pending",
	}

	return stage, nil
}

func generateID() string {
	return "stage-" + uuid.New().String()
}
