package docs

// @title           Taboo Game API
// @version         1.0
// @description     API Server for Taboo Game Application

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.basic  BasicAuth

// Game Management
// @Summary Create a new game
// @Description Create a new game session with specified team size
// @Tags games
// @Accept json
// @Produce json
// @Param teamSize body int true "Team size (3 or 4)"
// @Success 200 {object} models.Game
// @Router /games [post]

// @Summary Get game details
// @Description Get details of a specific game
// @Tags games
// @Accept json
// @Produce json
// @Param gameId path string true "Game ID"
// @Success 200 {object} models.Game
// @Router /games/{gameId} [get]

// Player Management
// @Summary Add player to game
// @Description Add a new player to an existing game
// @Tags players
// @Accept json
// @Produce json
// @Param gameId path string true "Game ID"
// @Param playerName body string true "Player name"
// @Success 200 {object} models.Player
// @Router /games/{gameId}/players [post]

// Match Management
// @Summary Start a match
// @Description Start a new match in a game
// @Tags matches
// @Accept json
// @Produce json
// @Param gameId path string true "Game ID"
// @Param matchId path string true "Match ID"
// @Param teamAssignments body object true "Team assignments"
// @Success 200 {object} models.MatchDetails
// @Router /games/{gameId}/matches/{matchId} [post]

// @Summary Process guess attempt
// @Description Process a player's guess attempt
// @Tags matches
// @Accept json
// @Produce json
// @Param matchId path string true "Match ID"
// @Param attempt body models.GuessAttempt true "Guess attempt details"
// @Success 200 {object} models.MatchDetails
// @Router /matches/{matchId}/guess [post]

// Team Management
// @Summary Switch player team
// @Description Move a player from one team to another
// @Tags teams
// @Accept json
// @Produce json
// @Param matchId path string true "Match ID"
// @Param playerId path string true "Player ID"
// @Success 200 {object} models.MatchDetails
// @Router /matches/{matchId}/teams/switch/{playerId} [post]
