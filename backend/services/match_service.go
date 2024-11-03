package services

import (
	"encoding/json"
	"errors"
	"taboo-game/models"
	"taboo-game/websocket"
	"time"

	"github.com/google/uuid"
)

type MatchService struct {
	matches      map[string]*models.MatchDetails
	words        []string
	gameService  GameServiceInterface
	wsManager    *websocket.Manager
	turnDuration time.Duration
}

func NewMatchService(gameService GameServiceInterface, wsManager *websocket.Manager) *MatchService {
	return &MatchService{
		matches:      make(map[string]*models.MatchDetails),
		words:        []string{},
		gameService:  gameService,
		wsManager:    wsManager,
		turnDuration: 60 * time.Second,
	}
}

func (s *MatchService) GetMatch(gameID, matchID string) (*models.MatchDetails, error) {
	match, exists := s.matches[matchID]
	if !exists {
		return nil, errors.New("match not found")
	}
	return match, nil
}

func (s *MatchService) StartMatch(gameID, matchID string, teamAssignments map[string][]string) (*models.MatchDetails, error) {
	// First verify the game exists
	_, err := s.gameService.GetGame(gameID)
	if err != nil {
		return nil, errors.New("game not found")
	}

	match, exists := s.matches[matchID]
	if !exists {
		match = &models.MatchDetails{
			ID:     matchID,
			GameID: gameID,
			Status: models.MatchStatusPending,
		}
		s.matches[matchID] = match
	}

	// Validate team assignments
	teamAPlayers, hasTeamA := teamAssignments["teamA"]
	teamBPlayers, hasTeamB := teamAssignments["teamB"]
	if !hasTeamA || !hasTeamB {
		return nil, errors.New("both teams must be assigned")
	}

	// Validate minimum team sizes
	minPlayersPerTeam := 2 // Minimum players needed for a valid team
	if len(teamAPlayers) < minPlayersPerTeam || len(teamBPlayers) < minPlayersPerTeam {
		return nil, errors.New("each team must have at least 2 players")
	}

	// Validate team balance (difference should not be more than 1 player)
	if abs(len(teamAPlayers)-len(teamBPlayers)) > 1 {
		return nil, errors.New("teams must be balanced (difference of max 1 player)")
	}

	// Update match with team assignments
	match.TeamAPlayers = teamAPlayers
	match.TeamBPlayers = teamBPlayers
	match.Status = models.MatchStatusPending
	match.CurrentWord = s.getNextWord()
	match.TeamATurn = true

	return match, nil
}

func (s *MatchService) ScorePoint(matchID string, isTeamA bool) (*models.MatchDetails, error) {
	match, exists := s.matches[matchID]
	if !exists {
		return nil, errors.New("match not found")
	}

	if match.Status != models.MatchStatusInProgress {
		return nil, errors.New("match is not in progress")
	}

	// Update score
	if isTeamA {
		match.TeamAScore++
	} else {
		match.TeamBScore++
	}

	// Broadcast score update
	scoreUpdateData := models.ScoreUpdateData{
		TeamAScore:  match.TeamAScore,
		TeamBScore:  match.TeamBScore,
		ScoringTeam: map[bool]string{true: "teamA", false: "teamB"}[isTeamA],
	}
	dataJSON, _ := json.Marshal(scoreUpdateData)

	scoreUpdate := models.GameEvent{
		Type:      models.EventTypeScoreUpdate,
		GameID:    match.GameID,
		MatchID:   match.ID,
		Timestamp: time.Now(),
		Data:      json.RawMessage(dataJSON),
	}

	eventJSON, _ := json.Marshal(scoreUpdate)
	s.wsManager.SendToGame(match.GameID, eventJSON)

	return match, nil
}

func (s *MatchService) ChangeTurn(matchID string) (*models.MatchDetails, error) {
	match, exists := s.matches[matchID]
	if !exists {
		return nil, errors.New("match not found")
	}

	// Switch turns
	match.TeamATurn = !match.TeamATurn
	match.CurrentWord = s.getNextWord()

	// Broadcast turn change
	turnChangeData := models.TurnChangeData{
		ActiveTeam: map[bool]string{true: "teamA", false: "teamB"}[match.TeamATurn],
		TimeLeft:   int(s.turnDuration.Seconds()),
	}
	dataJSON, _ := json.Marshal(turnChangeData)

	turnChange := models.GameEvent{
		Type:      models.EventTypeTurnChange,
		GameID:    match.GameID,
		MatchID:   match.ID,
		Timestamp: time.Now(),
		Data:      json.RawMessage(dataJSON),
	}

	eventJSON, _ := json.Marshal(turnChange)
	s.wsManager.SendToGame(match.GameID, eventJSON)

	// Start turn timer
	go s.startTurnTimer(match)

	return match, nil
}

func (s *MatchService) startTurnTimer(match *models.MatchDetails) {
	timer := time.NewTimer(s.turnDuration)
	<-timer.C

	// Time's up, change turns
	s.ChangeTurn(match.ID)
}

func (s *MatchService) EndMatch(gameID, matchID string) (*models.MatchDetails, error) {
	match, exists := s.matches[matchID]
	if !exists {
		return nil, errors.New("match not found")
	}

	match.Status = models.MatchStatusCompleted
	return match, nil
}

func (s *MatchService) getNextWord() string {
	return "placeholder"
}

func (s *MatchService) CreateStage(gameID, matchID string, details models.MatchStageDetails) (*models.MatchStage, error) {
	match, exists := s.matches[matchID]
	if !exists {
		return nil, errors.New("match not found")
	}

	if match.Status != models.MatchStatusPending {
		return nil, errors.New("match is not in pending state")
	}

	stage := &models.MatchStage{
		ID:             generateID(),
		MatchID:        matchID,
		ActiveTeamID:   details.ActiveTeamID,
		SpottingTeamID: details.SpottingTeamID,
		ClueGivers:     details.ClueGivers,
		Guessers:       details.Guessers,
		Spotters:       details.Spotters,
		Status:         "pending",
	}

	return stage, nil
}

func generateID() string {
	return "stage-" + uuid.New().String()
}

func (s *MatchService) SwitchTeam(matchID string, playerID string) (*models.MatchDetails, error) {
	match, exists := s.matches[matchID]
	if !exists {
		return nil, errors.New("match not found")
	}

	if match.Status != models.MatchStatusPending {
		return nil, errors.New("team switches are only allowed before match starts")
	}

	// Find and remove player from current team
	var currentTeam, otherTeam *[]string
	if containsPlayer(match.TeamAPlayers, playerID) {
		currentTeam = &match.TeamAPlayers
		otherTeam = &match.TeamBPlayers
	} else if containsPlayer(match.TeamBPlayers, playerID) {
		currentTeam = &match.TeamBPlayers
		otherTeam = &match.TeamAPlayers
	} else {
		return nil, errors.New("player not found in any team")
	}

	// Check if switch would create imbalance
	if len(*currentTeam)-1 < len(*otherTeam)-1 {
		return nil, errors.New("switch would create team imbalance")
	}

	// Perform the switch
	*currentTeam = removePlayer(*currentTeam, playerID)
	*otherTeam = append(*otherTeam, playerID)

	return match, nil
}

// Helper functions
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func containsPlayer(players []string, playerID string) bool {
	for _, p := range players {
		if p == playerID {
			return true
		}
	}
	return false
}

func removePlayer(players []string, playerID string) []string {
	result := make([]string, 0)
	for _, p := range players {
		if p != playerID {
			result = append(result, p)
		}
	}
	return result
}
