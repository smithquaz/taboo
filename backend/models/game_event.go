package models

import (
	"encoding/json"
	"time"
)

type GameEventType string

const (
	EventTypeScoreUpdate GameEventType = "score_update"
	EventTypeTurnChange  GameEventType = "turn_change"
	EventTypeGameEnd     GameEventType = "game_end"
)

type GameEvent struct {
	Type      GameEventType     `json:"type"`
	GameID    string           `json:"gameId"`
	MatchID   string           `json:"matchId"`
	Data      json.RawMessage  `json:"data"`
	Timestamp time.Time        `json:"timestamp"`
}

type ScoreUpdateData struct {
	TeamAScore int `json:"teamAScore"`
	TeamBScore int `json:"teamBScore"`
	ScoringTeam string `json:"scoringTeam"`
}

type TurnChangeData struct {
	ActiveTeam string `json:"activeTeam"`
	TimeLeft   int    `json:"timeLeft"`
} 