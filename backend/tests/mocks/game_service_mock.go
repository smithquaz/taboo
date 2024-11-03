package mocks

import (
	"taboo-game/models"
)

type MockGameService struct {
	// Mock implementation fields
	CreateGameFunc func(teamSize int) (*models.Game, error)
	AddPlayerFunc  func(gameID string, playerName string) (*models.Player, error)
	GetGameFunc    func(gameID string) (*models.Game, error)
	StartGameFunc  func(gameID string) (*models.Game, error)
	EndGameFunc    func(gameID string) (*models.Game, error)
	UpdateGameFunc func(game *models.Game) error
}

// Implement interface methods
func (m *MockGameService) CreateGame(teamSize int) (*models.Game, error) {
	return m.CreateGameFunc(teamSize)
}

func (m *MockGameService) AddPlayer(gameID string, playerName string) (*models.Player, error) {
	return m.AddPlayerFunc(gameID, playerName)
}

func (m *MockGameService) GetGame(gameID string) (*models.Game, error) {
	return m.GetGameFunc(gameID)
}

func (m *MockGameService) StartGame(gameID string) (*models.Game, error) {
	return m.StartGameFunc(gameID)
}

func (m *MockGameService) EndGame(gameID string) (*models.Game, error) {
	return m.EndGameFunc(gameID)
}

func (m *MockGameService) UpdateGame(game *models.Game) error {
	return m.UpdateGameFunc(game)
}
