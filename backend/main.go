package main

import (
	"log"
	"taboo-game/handlers"
	"taboo-game/helpers"
	"taboo-game/routes"
	"taboo-game/services"
	"taboo-game/websocket"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// CORS configuration
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:5173"} // Vite default port
	r.Use(cors.New(config))

	// Load words from CSV files
	commonWords, err := helpers.LoadWordsFromCSV("data/common_words.csv", "common")
	if err != nil {
		log.Fatalf("Failed to load common words: %v", err)
	}

	specificWords, err := helpers.LoadWordsFromCSV("data/domain_words.csv", "specific")
	if err != nil {
		log.Fatalf("Failed to load specific words: %v", err)
	}

	// Initialize WebSocket manager first
	wsManager := websocket.NewManager()
	go wsManager.Start()

	// Initialize services
	gameService := services.NewGameService()
	matchService := services.NewMatchService(gameService, wsManager)

	// Initialize handlers
	gameHandler := handlers.NewGameHandler(gameService)
	matchHandler := handlers.NewMatchHandler(matchService)

	// Initialize WebSocket handler
	wsHandler := handlers.NewWebSocketHandler(wsManager)

	// Initialize and register routes
	wordRoutes := routes.NewWordRoutes(commonWords, specificWords)
	gameRoutes := routes.NewGameRoutes(gameHandler)
	matchRoutes := routes.NewMatchRoutes(matchHandler)

	// Register all routes
	wordRoutes.RegisterRoutes(r)
	gameRoutes.RegisterRoutes(r)
	matchRoutes.RegisterRoutes(r)

	// Add WebSocket route
	r.GET("/ws/:gameId", wsHandler.HandleConnection)

	// Health check
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.Run(":8080")
}
