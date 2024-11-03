package main

import (
	"log"
	"taboo-game/docs"
	"taboo-game/handlers"
	"taboo-game/routes"
	"taboo-game/services"
	"taboo-game/websocket"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           Taboo Game API
// @version         1.0
// @description     API Server for Taboo Game Application
// @host            localhost:8080
// @BasePath        /api/v1
func main() {
	r := gin.Default()

	// CORS configuration
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:5173"} // Vite default port
	r.Use(cors.New(config))

	// Swagger documentation
	docs.SwaggerInfo.BasePath = "/api/v1"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Initialize core services
	gameService := services.NewGameService()
	wordService, err := services.NewWordService("data")
	if err != nil {
		log.Fatalf("Failed to initialize word service: %v", err)
	}

	// Initialize handlers first
	gameHandler := handlers.NewGameHandler(gameService)
	gameEventsHandler := handlers.NewGameEventsHandler(gameService)

	// Initialize websocket with the game events handler
	wsManager := websocket.NewManager(gameEventsHandler)
	go wsManager.Run()

	// Initialize services that depend on websocket
	matchService := services.NewMatchService(gameService, wsManager)
	services.NewGameEventsService(matchService, wordService, wsManager)

	// Initialize handlers that depend on services
	matchHandler := handlers.NewMatchHandler(matchService)

	// Register routes
	routes.SetupWebSocketRoutes(r, wsManager)
	routes.NewGameRoutes(gameHandler).RegisterRoutes(r)
	routes.NewMatchRoutes(matchHandler).RegisterRoutes(r)

	// Health check
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	// Start server
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
