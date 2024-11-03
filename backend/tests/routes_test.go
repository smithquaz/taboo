package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"taboo-game/handlers"
	"taboo-game/routes"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	// Initialize handlers
	teamHandler := handlers.NewTeamHandler()
	playerHandler := handlers.NewPlayerHandler()

	// Initialize routes
	teamRoutes := routes.NewTeamRoutes(teamHandler)
	teamRoutes.RegisterRoutes(router)

	playerRoutes := routes.NewPlayerRoutes(playerHandler)
	playerRoutes.RegisterRoutes(router)

	// Add the ping route for basic connectivity testing
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	return router
}

func TestRoutes(t *testing.T) {
	router := setupTestRouter()

	tests := []struct {
		name       string
		method     string
		path       string
		wantStatus int
	}{
		{"Ping", "GET", "/ping", 200},
		{"Create Team", "POST", "/api/v1/teams/", 200},
		{"Get Team", "GET", "/api/v1/teams/123", 200},
		{"Update Team", "PUT", "/api/v1/teams/123", 200},
		{"Update Team Players", "PUT", "/api/v1/teams/123/players", 200},
		{"Delete Team", "DELETE", "/api/v1/teams/123", 200},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(tt.method, tt.path, nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code, "Route %s %s returned wrong status code", tt.method, tt.path)
		})
	}
}

func TestPingRoute(t *testing.T) {
	router := setupTestRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, `{"message":"pong"}`, w.Body.String())
}
