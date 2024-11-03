package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
	"taboo-game/handlers"
	"taboo-game/helpers"
	"taboo-game/routes"
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

	// Initialize handlers
	gameHandler := handlers.NewGameHandler()
	playerHandler := handlers.NewPlayerHandler()
	teamHandler := handlers.NewTeamHandler()

	// Initialize and register routes
	wordRoutes := routes.NewWordRoutes(commonWords, specificWords)
	gameRoutes := routes.NewGameRoutes(gameHandler)
	playerRoutes := routes.NewPlayerRoutes(playerHandler)
	teamRoutes := routes.NewTeamRoutes(teamHandler)

	// Register all routes
	wordRoutes.RegisterRoutes(r)
	gameRoutes.RegisterRoutes(r)
	playerRoutes.RegisterRoutes(r)
	teamRoutes.RegisterRoutes(r)

	// Health check
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.Run(":8080")
}
