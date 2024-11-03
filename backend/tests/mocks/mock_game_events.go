package mocks

type MockGameEvents struct {
	StartStageFunc    func(gameID string, stageNum int) error
	HandleClueFunc    func(gameID string, playerID string, clue string) error
	HandleGuessFunc   func(gameID string, playerID string, guess string) error
	HandleViolationFunc func(gameID string, reporterID string, violationType string) error
}

func (m *MockGameEvents) StartStage(gameID string, stageNum int) error {
	if m.StartStageFunc != nil {
		return m.StartStageFunc(gameID, stageNum)
	}
	return nil
}

func (m *MockGameEvents) HandleClue(gameID string, playerID string, clue string) error {
	if m.HandleClueFunc != nil {
		return m.HandleClueFunc(gameID, playerID, clue)
	}
	return nil
}

func (m *MockGameEvents) HandleGuess(gameID string, playerID string, guess string) error {
	if m.HandleGuessFunc != nil {
		return m.HandleGuessFunc(gameID, playerID, guess)
	}
	return nil
}

func (m *MockGameEvents) HandleViolation(gameID string, reporterID string, violationType string) error {
	if m.HandleViolationFunc != nil {
		return m.HandleViolationFunc(gameID, reporterID, violationType)
	}
	return nil
} 