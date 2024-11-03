package tests

import (
	"path/filepath"
	"runtime"
	"testing"

	"taboo-game/helpers"
	"taboo-game/models"

	"github.com/stretchr/testify/assert"
)

// getTestDataPath returns the absolute path to the testdata directory
func getTestDataPath() string {
	_, b, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(b), "testdata")
}

func TestLoadWordsFromCSV(t *testing.T) {
	// Get path to test CSV file
	testFilePath := filepath.Join(getTestDataPath(), "words.csv")

	// Test loading words
	words, err := helpers.LoadWordsFromCSV(testFilePath, "animals")

	// Assert no error occurred
	assert.NoError(t, err)

	// Assert we got the expected number of words
	assert.Equal(t, 2, len(words))

	// Check first word
	assert.Equal(t, "dog", words[0].Word)
	assert.Equal(t, models.WordCategory("animals"), words[0].Category)
	assert.Equal(t, []string{"bark", "pet", "animal"}, words[0].TabooWords)
	assert.Equal(t, 1, words[0].Difficulty)

	// Check second word
	assert.Equal(t, "cat", words[1].Word)
	assert.Equal(t, models.WordCategory("animals"), words[1].Category)
	assert.Equal(t, []string{"meow", "kitten", "pet"}, words[1].TabooWords)
	assert.Equal(t, 1, words[1].Difficulty)

	// Test with non-existent file
	_, err = helpers.LoadWordsFromCSV("non_existent.csv", "animals")
	assert.Error(t, err)
}
