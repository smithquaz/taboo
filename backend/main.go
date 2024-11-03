package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
	"taboo-game/handlers"
	"taboo-game/helpers"
	"taboo-game/routes"
	"taboo-game/services"
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

	// Initialize services
	gameService := services.NewGameService()
	matchService := services.NewMatchService(gameService)

	// Initialize handlers
	gameHandler := handlers.NewGameHandler(gameService)
	matchHandler := handlers.NewMatchHandler(matchService)

	// Initialize and register routes
	wordRoutes := routes.NewWordRoutes(commonWords, specificWords)
	gameRoutes := routes.NewGameRoutes(gameHandler)
	matchRoutes := routes.NewMatchRoutes(matchHandler)

	// Register all routes
	wordRoutes.RegisterRoutes(r)
	gameRoutes.RegisterRoutes(r)
	matchRoutes.RegisterRoutes(r)

	// Health check
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.Run(":8080")
}
