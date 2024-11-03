package tests

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"taboo-game/handlers"
	"taboo-game/models"
	"taboo-game/tests/mocks"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCreateGame(t *testing.T) {
	// Setup
	mockGame := &models.Game{ID: "test-id", Status: "waiting"}
	mockService := &mocks.MockGameService{
		CreateGameFunc: func(teamSize int) (*models.Game, error) {
			if teamSize != 3 && teamSize != 4 {
				return nil, errors.New("invalid team size")
			}
			return mockGame, nil
		},
	}

	handler := handlers.NewGameHandler(mockService)

	// Test cases
	tests := []struct {
		name     string
		teamSize int
		want     *models.Game
		wantErr  bool
	}{
		{
			name:     "Valid team size",
			teamSize: 3,
			want:     mockGame,
			wantErr:  false,
		},
		{
			name:     "Invalid team size",
			teamSize: 5,
			want:     nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create request
			req := httptest.NewRequest("POST", "/games",
				strings.NewReader(fmt.Sprintf(`{"teamSize":%d}`, tt.teamSize)))
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = req

			// Call handler
			handler.CreateGame(c)

			// Assert response
			if tt.wantErr {
				assert.Equal(t, http.StatusBadRequest, w.Code)
			} else {
				assert.Equal(t, http.StatusCreated, w.Code)
				var response models.Game
				err := json.NewDecoder(w.Body).Decode(&response)
				assert.NoError(t, err)
				assert.Equal(t, mockGame.ID, response.ID)
			}
		})
	}
}
