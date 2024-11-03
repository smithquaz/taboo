package services

import (
	"errors"
	"github.com/google/uuid"
	"taboo-game/models"
	"time"
)

type GameService struct {
	games map[string]*models.Game
}

func NewGameService() *GameService {
	return &GameService{
		games: make(map[string]*models.Game),
	}
}

func (s *GameService) CreateGame(teamSize int) (*models.Game, error) {
	game := &models.Game{
		ID:        uuid.New().String(),
		CreatedAt: time.Now(),
		Status:    "waiting",
		Teams: []models.Team{
			{
				ID:      uuid.New().String(),
				Name:    "Team 1",
				Players: []models.Player{},
				Score:   0,
				Size:    teamSize,
			},
			{
				ID:      uuid.New().String(),
				Name:    "Team 2",
				Players: []models.Player{},
				Score:   0,
				Size:    teamSize,
			},
		},
		Matches: []models.Match{},
	}

	s.games[game.ID] = game
	return game, nil
}

func (s *GameService) AddPlayer(gameID string, playerName string) (*models.Player, error) {
	game, exists := s.games[gameID]
	if !exists {
		return nil, errors.New("game not found")
	}

	if game.Status != "waiting" {
		return nil, errors.New("game has already started")
	}

	// Find team with space available
	var targetTeam *models.Team
	for i := range game.Teams {
		if len(game.Teams[i].Players) < game.Teams[i].Size {
			targetTeam = &game.Teams[i]
			break
		}
	}

	if targetTeam == nil {
		return nil, errors.New("game is full")
	}

	player := models.Player{
		ID:       uuid.New().String(),
		Name:     playerName,
		TeamID:   targetTeam.ID,
		JoinedAt: time.Now(),
	}

	targetTeam.Players = append(targetTeam.Players, player)
	return &player, nil
}
