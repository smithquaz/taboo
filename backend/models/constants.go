package models

// GameStatus represents the possible states of a game
type GameStatus string

const (
	GameStatusWaiting    GameStatus = "waiting"
	GameStatusInProgress GameStatus = "in_progress"
	GameStatusCompleted  GameStatus = "completed"
)

// StageStatus represents the possible states of a stage
type StageStatus string

const (
	StageStatusPending   StageStatus = "pending"
	StageStatusActive    StageStatus = "active"
	StageStatusCompleted StageStatus = "completed"
)

// WordCategory represents the possible categories for words
type WordCategory string

const (
	WordCategoryCommon WordCategory = "common"
	WordCategoryDomain WordCategory = "domain"
)
