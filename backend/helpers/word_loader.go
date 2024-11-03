package helpers

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"

	"taboo-game/models"

	"github.com/google/uuid"
)

func LoadWordsFromCSV(filename string, category models.WordCategory) ([]models.Word, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("error reading CSV: %v", err)
	}

	var words []models.Word
	for i, record := range records {
		if i == 0 { // Skip header row
			continue
		}
		if len(record) < 2 {
			continue // Skip invalid rows
		}

		word := models.Word{
			ID:         uuid.New().String(),
			Word:       strings.TrimSpace(record[0]),
			TabooWords: make([]string, 0),
			Category:   category,
			Difficulty: 1, // Default difficulty, adjust as needed
		}

		// Add taboo words from remaining columns
		for j := 1; j < len(record); j++ {
			if tabooWord := strings.TrimSpace(record[j]); tabooWord != "" {
				word.TabooWords = append(word.TabooWords, tabooWord)
			}
		}

		words = append(words, word)
	}

	return words, nil
}
