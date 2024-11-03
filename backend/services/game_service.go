package services

import (
	"errors"
	"taboo-game/models"
	"time"

	"github.com/google/uuid"
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

func (s *GameService) GetGame(gameID string) (*models.Game, error) {
	game, exists := s.games[gameID]
	if !exists {
		return nil, errors.New("game not found")
	}
	return game, nil
}

func (s *GameService) StartGame(gameID string) (*models.Game, error) {
	game, exists := s.games[gameID]
	if !exists {
		return nil, errors.New("game not found")
	}

	// Validate team sizes
	for _, team := range game.Teams {
		if len(team.Players) != team.Size {
			return nil, errors.New("teams must be full to start game")
		}
	}

	game.Status = models.GameStatusInProgress
	game.Matches = []models.Match{
		createMatch(1, gameID),
		createMatch(2, gameID),
		createMatch(3, gameID),
	}

	return game, nil
}

func (s *GameService) EndGame(gameID string) (*models.Game, error) {
	game, exists := s.games[gameID]
	if !exists {
		return nil, errors.New("game not found")
	}

	if game.Status != models.GameStatusInProgress {
		return nil, errors.New("game is not in progress")
	}

	game.Status = models.GameStatusCompleted
	return game, nil
}

// Helper function to create a match
func createMatch(number int, gameID string) models.Match {
	return models.Match{
		ID:     uuid.New().String(),
		GameID: gameID,
		Number: number,
		Status: models.MatchStatusPending,
		Stages: make([]models.Stage, 0),
	}
}

func (s *GameService) UpdateGame(game *models.Game) error {
	if _, exists := s.games[game.ID]; !exists {
		return errors.New("game not found")
	}
	s.games[game.ID] = game
	return nil
}
