package services

import (
	"encoding/csv"
	"math/rand"
	"os"
	"time"
)

type WordService struct {
	words     []string
	usedWords map[string]bool
}

func NewWordService() *WordService {
	return &WordService{
		words:     make([]string, 0),
		usedWords: make(map[string]bool),
	}
}

func (s *WordService) LoadWordsFromCSV(filepath string) error {
	file, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	for _, record := range records {
		if len(record) > 0 {
			s.words = append(s.words, record[0])
		}
	}

	return nil
}

func (s *WordService) GetNextWord() string {
	// If all words have been used, reset the used words
	if len(s.usedWords) >= len(s.words) {
		s.usedWords = make(map[string]bool)
	}

	// Get a random unused word
	rand.Seed(time.Now().UnixNano())
	availableWords := make([]string, 0)
	for _, word := range s.words {
		if !s.usedWords[word] {
			availableWords = append(availableWords, word)
		}
	}

	if len(availableWords) == 0 {
		return ""
	}

	selectedWord := availableWords[rand.Intn(len(availableWords))]
	s.usedWords[selectedWord] = true
	return selectedWord
}
