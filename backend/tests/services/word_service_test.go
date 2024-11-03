package services_test

import (
	"encoding/csv"
	"os"
	"path/filepath"
	"taboo-game/services"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWordService(t *testing.T) {
	// Create test data directory and files
	testDir := "testdata"
	err := os.MkdirAll(testDir, 0755)
	assert.NoError(t, err)
	defer os.RemoveAll(testDir)

	// Create test word files
	testWords := [][]string{
		{"word", "taboo1", "taboo2", "taboo3", "1", "general"},
		{"test", "no", "yes", "maybe", "2", "test"},
	}

	files := []string{"common_words.csv", "domain_words.csv"}
	for _, file := range files {
		f, err := os.Create(filepath.Join(testDir, file))
		assert.NoError(t, err)

		w := csv.NewWriter(f)
		w.Write([]string{"target", "taboo1", "taboo2", "taboo3", "difficulty", "category"})
		w.WriteAll(testWords)
		f.Close()
	}

	// Initialize service with test directory
	ws, err := services.NewWordService(testDir)
	assert.NoError(t, err)
	assert.NotNil(t, ws)

	// Test getting unique cards
	usedCards := make(map[string]bool)
	expectedCards := len(testWords) * len(files)

	// Add debug logging
	t.Logf("Expecting to get %d cards", expectedCards)

	for i := 0; i < expectedCards; i++ {
		card, err := ws.GetNextCard()
		if err != nil {
			t.Logf("Failed to get card #%d: %v", i+1, err)
			t.FailNow()
		}
		if card == nil {
			t.Logf("Got nil card on iteration #%d", i+1)
			t.FailNow()
		}

		t.Logf("Got card #%d: %s", i+1, card.ID)
		assert.False(t, usedCards[card.ID], "Card should not be repeated")
		usedCards[card.ID] = true
	}

	// Verify we got all expected cards
	assert.Equal(t, expectedCards, len(usedCards), "Should have received all unique cards")
}
