package services

import (
	"encoding/csv"
	"errors"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"sync"

	"taboo-game/models"
)

type WordService struct {
	wordCards []models.WordCard
	mu        sync.RWMutex
	usedCards map[string]bool
}

func NewWordService(dataDir string) (*WordService, error) {
	ws := &WordService{
		usedCards: make(map[string]bool),
	}

	if err := ws.loadWords(dataDir); err != nil {
		return nil, err
	}

	return ws, nil
}

func (ws *WordService) loadWords(dataDir string) error {
	files := []string{
		filepath.Join(dataDir, "common_words.csv"),
		filepath.Join(dataDir, "domain_words.csv"),
	}

	for _, file := range files {
		if err := ws.loadWordsFromFile(file); err != nil {
			return err
		}
	}

	return nil
}

func (ws *WordService) loadWordsFromFile(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)

	// Skip header row
	_, err = reader.Read()
	if err != nil {
		return err
	}

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		// Assuming CSV format: target_word,taboo_word1,taboo_word2,taboo_word3,difficulty,category
		if len(record) < 6 {
			continue
		}

		tabooWords := record[1:4]
		card := models.WordCard{
			ID:         filepath.Base(filePath) + "-" + record[0],
			TargetWord: record[0],
			TabooWords: tabooWords,
			Difficulty: 1, // Parse from record[4] if needed
			Category:   record[5],
		}

		ws.wordCards = append(ws.wordCards, card)
	}

	return nil
}

func (ws *WordService) GetNextCard() (*models.WordCard, error) {
	ws.mu.Lock()
	defer ws.mu.Unlock()

	// Reset used cards if all have been used
	if len(ws.usedCards) == len(ws.wordCards) {
		ws.usedCards = make(map[string]bool)
	}

	// Find an unused card
	for attempts := 0; attempts < 10; attempts++ {
		idx := rand.Intn(len(ws.wordCards))
		card := ws.wordCards[idx]

		if !ws.usedCards[card.ID] {
			ws.usedCards[card.ID] = true
			return &card, nil
		}
	}

	return nil, errors.New("failed to find unused card")
}
