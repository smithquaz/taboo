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
	GetMatch(gameID, matchID string) (*models.MatchDetails, error)
	StartMatch(gameID, matchID string, teamAssignments map[string][]string) (*models.MatchDetails, error)
	EndMatch(gameID, matchID string) (*models.MatchDetails, error)
	ScorePoint(matchID string, isTeamA bool) (*models.MatchDetails, error)
	CreateStage(gameID, matchID string, stageDetails models.MatchStageDetails) (*models.MatchStage, error)
	SwitchTeam(matchID string, playerID string) (*models.MatchDetails, error)
	ProcessGuessAttempt(gameID, matchID string, attempt *models.GuessAttempt) error
}
