package models

type MatchStatus string

const (
	MatchStatusPending    MatchStatus = "pending"
	MatchStatusInProgress MatchStatus = "in_progress"
	MatchStatusCompleted  MatchStatus = "completed"
)

type MatchDetails struct {
	ID           string      `json:"id"`
	GameID       string      `json:"gameId"`
	Status       MatchStatus `json:"status"`
	TeamATurn    bool        `json:"teamATurn"`
	TeamAScore   int         `json:"teamAScore"`
	TeamBScore   int         `json:"teamBScore"`
	TeamAPlayers []string    `json:"teamAPlayers"`
	TeamBPlayers []string    `json:"teamBPlayers"`
	CurrentWord  string      `json:"currentWord"`
}

type MatchStage struct {
	ID             string   `json:"id"`
	MatchID        string   `json:"matchId"`
	ActiveTeamID   string   `json:"activeTeamId"`
	SpottingTeamID string   `json:"spottingTeamId"`
	ClueGivers     []string `json:"clueGivers"`
	Guessers       []string `json:"guessers"`
	Spotters       []string `json:"spotters"`
	Status         string   `json:"status"`
}

type MatchStageDetails struct {
	ActiveTeamID   string   `json:"activeTeamId"`
	SpottingTeamID string   `json:"spottingTeamId"`
	ClueGivers     []string `json:"clueGivers"`
	Guessers       []string `json:"guessers"`
	Spotters       []string `json:"spotters"`
}
