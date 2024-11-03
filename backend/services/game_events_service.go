package services

import (
	"encoding/json"
	"sync"
	"taboo-game/websocket"
	"time"
)

type GameEventsService struct {
	matchService *MatchService
	wordService  *WordService
	wsManager    *websocket.Manager
	activeStages map[string]*StageTimer
	mu           sync.RWMutex
}

type StageTimer struct {
	timer   *time.Timer
	ticker  *time.Ticker
	done    chan bool
	endTime time.Time
}

func NewGameEventsService(ms *MatchService, ws *WordService, wm *websocket.Manager) *GameEventsService {
	return &GameEventsService{
		matchService: ms,
		wordService:  ws,
		wsManager:    wm,
		activeStages: make(map[string]*StageTimer),
	}
}

func (s *GameEventsService) StartStage(gameID string, stageNum int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Get new word card
	wordCard, err := s.wordService.GetNextCard()
	if err != nil {
		return err
	}

	duration := 3 * time.Minute
	// Initialize stage timer
	timer := &StageTimer{
		timer:   time.NewTimer(duration),
		ticker:  time.NewTicker(time.Second),
		done:    make(chan bool),
		endTime: time.Now().Add(duration),
	}
	s.activeStages[gameID] = timer

	// Start timer goroutine
	go s.runStageTimer(gameID, timer)

	// Broadcast stage start
	msg := websocket.Message{
		Type:   websocket.StartStage,
		GameID: gameID,
		Payload: map[string]interface{}{
			"stage_num": stageNum,
			"word_card": wordCard,
			"duration":  int(duration.Seconds()),
		},
	}
	data, _ := json.Marshal(msg)
	s.wsManager.SendToGame(gameID, data)

	return nil
}

func (s *GameEventsService) runStageTimer(gameID string, timer *StageTimer) {
	defer timer.ticker.Stop()
	defer timer.timer.Stop()

	for {
		select {
		case <-timer.ticker.C:
			remaining := time.Until(timer.endTime).Seconds()
			msg := websocket.Message{
				Type:   websocket.TimerUpdate,
				GameID: gameID,
				Payload: map[string]interface{}{
					"remaining": int(remaining),
				},
			}
			data, _ := json.Marshal(msg)
			s.wsManager.SendToGame(gameID, data)
		case <-timer.timer.C:
			s.handleStageEnd(gameID)
			return
		case <-timer.done:
			return
		}
	}
}

func (s *GameEventsService) handleStageEnd(gameID string) {
	// Update match state
	match, err := s.matchService.GetCurrentMatch(gameID)
	if err != nil {
		return
	}

	// Move to next stage or end match
	// Hard coded to 4 stages for now
	if match.CurrentStage.Number < 4 {
		s.matchService.NextStage(gameID)
	} else {
		s.matchService.EndCurrentMatch(gameID)
	}
}

func (s *GameEventsService) HandleClue(gameID, playerID, clue string) error {
	// Validate clue giver's turn
	// Broadcast clue to all players
	return nil
}

func (s *GameEventsService) HandleGuess(gameID, playerID, guess string) error {
	// Validate guess and update score
	// Broadcast result to all players
	return nil
}

func (s *GameEventsService) HandleViolation(gameID, reporterID, violationType string) error {
	// Validate violation report and update score
	// Broadcast violation to all players
	return nil
}
func encodeMessage(msg websocket.Message) []byte {
	data, _ := json.Marshal(msg)
	return data
}

func (s *GameEventsService) ProcessEvent(eventData []byte) error {
	var event websocket.Message
	if err := json.Unmarshal(eventData, &event); err != nil {
		return err
	}
	return s.HandleGameEvent(event)
}

func (s *GameEventsService) HandleGameEvent(event websocket.Message) error {
	// Implement game event handling logic
	return nil
}
