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

	// Find current team of the player
	isInTeamA := containsPlayer(match.TeamAPlayers, playerID)
	isInTeamB := containsPlayer(match.TeamBPlayers, playerID)

	if !isInTeamA && !isInTeamB {
		return nil, errors.New("player not found in any team")
	}

	// Calculate new team sizes after switch
	var newTeamASize, newTeamBSize int
	if isInTeamA {
		newTeamASize = len(match.TeamAPlayers) - 1
		newTeamBSize = len(match.TeamBPlayers) + 1
	} else {
		newTeamASize = len(match.TeamAPlayers) + 1
		newTeamBSize = len(match.TeamBPlayers) - 1
	}

	// Check if switch would create imbalance (difference > 1)
	if abs(newTeamASize-newTeamBSize) > 1 {
		return nil, errors.New("switch would create team imbalance")
	}

	// Perform the switch
	if isInTeamA {
		match.TeamAPlayers = removePlayer(match.TeamAPlayers, playerID)
		match.TeamBPlayers = append(match.TeamBPlayers, playerID)
	} else {
		match.TeamBPlayers = removePlayer(match.TeamBPlayers, playerID)
		match.TeamAPlayers = append(match.TeamAPlayers, playerID)
	}

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

func (s *MatchService) GetCurrentMatch(stageID string) (*models.MatchDetails, error) {
	// Find match containing this stage
	for _, match := range s.matches {
		if match.CurrentStage != nil && match.CurrentStage.ID == stageID {
			return match, nil
		}
	}
	return nil, errors.New("match not found for stage")
}

func (s *MatchService) getOpposingTeamID(teamID string) string {
	// Using TeamA/TeamB style
	if teamID == "teamA" {
		return "teamB"
	}
	return "teamA"
}

func (s *MatchService) getSmallerTeam(match *models.MatchDetails) string {
	// Using TeamA/TeamB style
	if len(match.TeamAPlayers) < len(match.TeamBPlayers) {
		return "teamA"
	}
	return "teamB"
}

func (s *MatchService) ProcessGuessAttempt(gameID, matchID string, attempt *models.GuessAttempt) error {
	match, exists := s.matches[matchID]
	if !exists || match.GameID != gameID {
		return errors.New("match not found for the given game")
	}

	if match.CurrentStage == nil {
		return errors.New("no active stage")
	}

	// Update both match and stage scores
	if attempt.Correct {
		if attempt.TeamID == "teamA" {
			match.TeamAScore += models.PointsCorrectGuess
			match.CurrentStage.TeamAScore += models.PointsCorrectGuess
		} else {
			match.TeamBScore += models.PointsCorrectGuess
			match.CurrentStage.TeamBScore += models.PointsCorrectGuess
		}
	}

	if attempt.Violation {
		opposingTeamID := s.getOpposingTeamID(attempt.TeamID)
		if opposingTeamID == "teamA" {
			match.TeamAScore += models.PointsViolationCatch
			match.CurrentStage.TeamAScore += models.PointsViolationCatch
		} else {
			match.TeamBScore += models.PointsViolationCatch
			match.CurrentStage.TeamBScore += models.PointsViolationCatch
		}
	}

	// Emit score update event
	dataJSON, _ := json.Marshal(struct {
		TeamAScore int `json:"teamAScore"`
		TeamBScore int `json:"teamBScore"`
	}{
		TeamAScore: match.TeamAScore,
		TeamBScore: match.TeamBScore,
	})

	scoreUpdate := models.GameEvent{
		Type:    "SCORE_UPDATE",
		GameID:  match.GameID,
		MatchID: match.ID,
		Data:    json.RawMessage(dataJSON),
	}

	eventJSON, _ := json.Marshal(scoreUpdate)
	s.wsManager.SendToGame(match.GameID, eventJSON)

	return nil
}

func (s *MatchService) FinalizeStageScores(stageID string) error {
	match, err := s.GetCurrentMatch(stageID)
	if err != nil {
		return err
	}

	// Apply team size balance adjustment at the end of each stage
	smallerTeam := s.getSmallerTeam(match)
	if smallerTeam == "teamA" && len(match.TeamAPlayers) == 3 {
		match.TeamAScore += models.BasePointsTeamOfThree
		match.CurrentStage.TeamAScore += models.BasePointsTeamOfThree
	} else if smallerTeam == "teamB" && len(match.TeamBPlayers) == 3 {
		match.TeamBScore += models.BasePointsTeamOfThree
		match.CurrentStage.TeamBScore += models.BasePointsTeamOfThree
	}

	return nil
}

// Add this method to store matches for testing
func (s *MatchService) StoreMatch(match *models.MatchDetails) {
	s.matches[match.ID] = match
}
