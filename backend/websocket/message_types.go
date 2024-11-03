package websocket

// Message represents a WebSocket message structure
type Message struct {
	Type     MessageType            `json:"type"`
	GameID   string                 `json:"gameId"`
	PlayerID string                 `json:"playerId"`
	Payload  map[string]interface{} `json:"payload"`
}

type MessageType string

const (
	StartStage  MessageType = "START_STAGE"
	GiveClue    MessageType = "GIVE_CLUE"
	TimerUpdate MessageType = "TIMER_UPDATE"
	StageEnd    MessageType = "STAGE_END"
	GameEnd     MessageType = "GAME_END"
)
