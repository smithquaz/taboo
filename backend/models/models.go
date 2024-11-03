package models

import "time"

// Player represents a user in the game
type Player struct {
	ID       string    `json:"id"`
	Name     string    `json:"name"`
	TeamID   string    `json:"teamId"`
	JoinedAt time.Time `json:"joinedAt"`
}

// Team represents a group of players
type Team struct {
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	GameID  string   `json:"gameId"`
	Players []Player `json:"players"`
	Score   int      `json:"score"`
	Size    int      `json:"size"` // 3 or 4 players
}

// Game represents an entire game session
type Game struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	Status    string    `json:"status"` // "waiting", "in_progress", "completed"
	Teams     []Team    `json:"teams"`
	Matches   []Match   `json:"matches"`
}

// Match represents one of the three matches in a game
type Match struct {
	ID        string    `json:"id"`
	GameID    string    `json:"gameId"`
	Number    int       `json:"number"` // 1, 2, or 3
	Status    string    `json:"status"` // "pending", "active", "completed"
	Stages    []Stage   `json:"stages"`
	StartedAt time.Time `json:"startedAt,omitempty"`
	EndedAt   time.Time `json:"endedAt,omitempty"`
}

// Stage represents one of the four stages within a match
type Stage struct {
	ID      string `json:"id"`
	MatchID string `json:"matchId"`
	Number  int    `json:"number"` // 1, 2, 3, or 4

	// Active team details
	ActiveTeamID string   `json:"activeTeamId"`
	ClueGivers   []string `json:"clueGivers"` // Player IDs of the 2 clue givers
	Guessers     []string `json:"guessers"`   // Player IDs of the guessers

	// Spotting team details
	SpottingTeamID string   `json:"spottingTeamId"`
	Spotters       []string `json:"spotters"` // Player IDs of the 2 spotters

	Score     int       `json:"score"`
	Duration  int       `json:"duration"` // 180 seconds (3 minutes)
	Status    string    `json:"status"`   // "pending", "active", "completed"
	StartedAt time.Time `json:"startedAt,omitempty"`
	EndedAt   time.Time `json:"endedAt,omitempty"`
}

// Word represents a word to be guessed and its taboo words
type Word struct {
	ID          string   `json:"id"`
	Word        string   `json:"word"`
	TabooWords  []string `json:"tabooWords"`
	Category    string   `json:"category"` // "common" or "domain"
	Difficulty  int      `json:"difficulty"` // 1-5
}
