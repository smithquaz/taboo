package services

import (
	"errors"
	"fmt"
	"taboo-game/models"
	"time"
)

type MatchService struct {
	gameService *GameService
}

func NewMatchService(gameService *GameService) *MatchService {
	return &MatchService{
		gameService: gameService,
	}
}

func (s *MatchService) GetMatch(gameID, matchID string) (*models.Match, error) {
	game, err := s.gameService.GetGame(gameID)
	if err != nil {
		return nil, err
	}

	for _, match := range game.Matches {
		if match.ID == matchID {
			return &match, nil
		}
	}
	return nil, errors.New("match not found")
}

func (s *MatchService) StartMatch(gameID, matchID string, teamAssignments map[string][]string) (*models.Match, error) {
	game, err := s.gameService.GetGame(gameID)
	if err != nil {
		return nil, err
	}

	var match *models.Match
	for i := range game.Matches {
		if game.Matches[i].ID == matchID {
			match = &game.Matches[i]
			break
		}
	}
	if match == nil {
		return nil, errors.New("match not found")
	}

	// Validate team assignments
	if err := s.validateTeamAssignments(game, teamAssignments); err != nil {
		return nil, err
	}

	// Update team assignments
	for teamID, playerIDs := range teamAssignments {
		for i := range game.Teams {
			if game.Teams[i].ID == teamID {
				game.Teams[i].Players = make([]models.Player, 0)
				for _, playerID := range playerIDs {
					// Add player to team
					// Note: In a real implementation, you'd want to validate these players exist
					game.Teams[i].Players = append(game.Teams[i].Players, models.Player{ID: playerID})
				}
			}
		}
	}

	match.Status = models.MatchStatusActive
	match.StartedAt = time.Now()

	// Save the updated game state
	if err := s.gameService.UpdateGame(game); err != nil {
		return nil, err
	}

	return match, nil
}

func (s *MatchService) EndMatch(gameID, matchID string) (*models.Match, error) {
	game, err := s.gameService.GetGame(gameID)
	if err != nil {
		return nil, err
	}

	var match *models.Match
	for i := range game.Matches {
		if game.Matches[i].ID == matchID {
			match = &game.Matches[i]
			break
		}
	}
	if match == nil {
		return nil, errors.New("match not found")
	}

	if match.Status != models.MatchStatusActive {
		return nil, errors.New("match is not active")
	}

	match.Status = models.MatchStatusCompleted
	match.EndedAt = time.Now()

	// Save the updated game state
	if err := s.gameService.UpdateGame(game); err != nil {
		return nil, err
	}

	return match, nil
}

func (s *MatchService) validateTeamAssignments(game *models.Game, assignments map[string][]string) error {
	// Validate that all teams are assigned the correct number of players
	for teamID, playerIDs := range assignments {
		var team *models.Team
		for i := range game.Teams {
			if game.Teams[i].ID == teamID {
				team = &game.Teams[i]
				break
			}
		}
		if team == nil {
			return fmt.Errorf("team %s not found", teamID)
		}

		if len(playerIDs) != team.Size {
			return fmt.Errorf("team %s requires %d players, got %d", teamID, team.Size, len(playerIDs))
		}
	}
	return nil
}
