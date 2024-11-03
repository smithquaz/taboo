package services

import "taboo-game/models"

type GameServiceInterface interface {
	CreateGame(teamSize int) (*models.Game, error)
	AddPlayer(gameID string, playerName string) (*models.Player, error)
	GetGame(gameID string) (*models.Game, error)
	StartGame(gameID string) (*models.Game, error)
	EndGame(gameID string) (*models.Game, error)
	UpdateGame(game *models.Game) error
}

type MatchServiceInterface interface {
	GetMatch(gameID, matchID string) (*models.Match, error)
	StartMatch(gameID, matchID string, teamAssignments map[string][]string) (*models.Match, error)
	EndMatch(gameID, matchID string) (*models.Match, error)
}
