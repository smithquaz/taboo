package models

type WordCard struct {
	ID         string   `json:"id"`
	TargetWord string   `json:"targetWord"`
	TabooWords []string `json:"tabooWords"`
	Difficulty int      `json:"difficulty"` // 1-3
	Category   string   `json:"category"`
}

type GuessAttempt struct {
	CardID      string `json:"cardId"`
	Correct     bool   `json:"correct"`
	Violation   bool   `json:"violation"`
	TeamID      string `json:"teamId"`
	StageID     string `json:"stageId"`
	TimestampMS int64  `json:"timestampMs"`
}
