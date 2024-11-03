# Taboo Backend

## Architecture Overview

### 1. Core Components

#### Service Layer
- **GameService**: Manages game lifecycle (creation, player management, state)
- **MatchService**: Handles match logic (stages, scoring, team management)
- **WordService**: Manages word cards and selection
- **GameEventsService**: Orchestrates real-time game events and timers
- **TeamService**: Manages team composition and balancing
- **PlayerService**: Handles player state and roles

#### Handler Layer
- **GameHandler**: HTTP endpoints for game management
- **MatchHandler**: HTTP endpoints for match operations
- **WebSocketHandler**: WebSocket connection management
- **TeamHandler**: Team management endpoints
- **PlayerHandler**: Player management endpoints

#### WebSocket Layer
- **WebSocket Manager**: Central hub for real-time communication
- **Client**: Represents a connected player
- **Message Types**: Standardized communication protocol

### 2. Dependency Flow

```
main.go
├─> WebSocketManager (real-time communication)
├─> GameService (game state)
├─> WordService (word management)
├─> TeamService (team management)
├─> PlayerService (player management)
├─> MatchService (match logic)
└─> GameEventsService (real-time events)
    ├─> MatchService
    ├─> WordService
    └─> WebSocketManager
```

### 3. API Endpoints

#### Player Management
```
POST   /api/v1/players          # Create new player
GET    /api/v1/players/:id      # Get player details
PUT    /api/v1/players/:id      # Update player
DELETE /api/v1/players/:id      # Remove player from game
```

#### Team Management
```
POST   /api/v1/games/:id/teams          # Create team
GET    /api/v1/games/:id/teams          # List teams
PUT    /api/v1/games/:id/teams/:teamId  # Update team
POST   /api/v1/teams/:id/players        # Add player to team
DELETE /api/v1/teams/:id/players/:pid   # Remove player from team
```

### 4. Real-time Communication Flow

1. **Client Connection**
```go
// Client connects through WebSocket endpoint
GET /ws/:game_id/:player_id
```

2. **Message Flow**
```
Client -> WebSocket -> GameEventsService -> Game Logic -> Broadcast to Clients
```

3. **Event Types**
```go
const (
    JoinGame        = "join_game"
    StartMatch      = "start_match"
    GiveClue        = "give_clue"
    MakeGuess       = "make_guess"
    ReportViolation = "report_violation"
    TimerUpdate     = "timer_update"
    // ...
)
```

### 5. Interface-based Architecture

```go
// Core service interfaces
type GameEventsServiceInterface interface {
    StartStage(gameID string, stageNum int) error
    HandleClue(gameID string, playerID string, clue string) error
    HandleGuess(gameID string, playerID string, guess string) error
    HandleViolation(gameID string, reporterID string, violationType string) error
}

type WebSocketManagerInterface interface {
    Register(client WebSocketClientInterface)
    SendToGame(gameID string, message []byte)
    HandleConnection(w http.ResponseWriter, r *http.Request, gameID, playerID string)
    Run()
}
```

### 6. Testing

The architecture supports easy testing through mock implementations:
```go
// Mock example
type MockWebSocketManager struct {
    SendToGameFunc func(gameID string, message []byte)
    RegisterFunc   func(client types.WebSocketClientInterface)
}
```

## Running the Server

### Local Development
```bash
go run main.go
```

### Docker

#### Build
```bash
docker build --no-cache --platform=linux/amd64 -t taboo:[version] .
```

#### Run
```bash
docker run -p 8080:8080 taboo:[version]
```

## API Documentation

Swagger documentation available at:
```
http://localhost:8080/swagger/*
```

## WebSocket Events

### Game Stage Flow
1. **Stage Start**
   - Server sends word card
   - Starts 3-minute timer
   - Broadcasts stage status

2. **During Stage**
   - Clue giving/guessing
   - Violation reporting
   - Score updates
   - Timer updates

3. **Stage End**
   - Final scoring
   - Next stage preparation
   - Team role rotation

### Message Format
```json
{
    "type": "message_type",
    "game_id": "game_identifier",
    "player_id": "player_identifier",
    "payload": {
        // Event-specific data
    }
}
```

## Development Notes

1. **Dependency Injection**
   - Services receive dependencies via constructors
   - Enables easy testing and implementation swapping

2. **Concurrency**
   - WebSocket manager runs in separate goroutine
   - Each client connection has Read/Write goroutines
   - Timer management for game stages

3. **Error Handling**
   - Service-level validation
   - WebSocket connection management
   - Game state consistency checks

## Component Details

### Player Handler
```go
type PlayerHandler struct {
    playerService PlayerServiceInterface
}

// Key responsibilities:
// - Player creation and validation
// - Player state management
// - Role assignment
// - Session management
```

### Team Handler
```go
type TeamHandler struct {
    teamService TeamServiceInterface
}

// Key responsibilities:
// - Team creation and validation
// - Team size balancing (3v4)
// - Player assignment to teams
// - Team role rotation
```

### Player Service
```go
type PlayerService interface {
    CreatePlayer(name string) (*models.Player, error)
    AssignRole(playerID string, role models.PlayerRole) error
    UpdatePlayerState(playerID string, state models.PlayerState) error
    RemovePlayer(playerID string) error
}
```

### Team Service
```go
type TeamService interface {
    CreateTeam(gameID string, size int) (*models.Team, error)
    AddPlayer(teamID string, playerID string) error
    RemovePlayer(teamID string, playerID string) error
    RotateRoles(teamID string) error
    ValidateTeamBalance(teamA, teamB *models.Team) error
}
```

## Data Models

### Player Model
```go
type Player struct {
    ID        string      `json:"id"`
    Name      string      `json:"name"`
    Role      PlayerRole  `json:"role"`
    TeamID    string      `json:"team_id"`
    State     PlayerState `json:"state"`
    Connected bool        `json:"connected"`
}
```

### Team Model
```go
type Team struct {
    ID       string   `json:"id"`
    GameID   string   `json:"game_id"`
    Players  []string `json:"players"`
    Size     int      `json:"size"`
    Score    int      `json:"score"`
    IsTeamA  bool     `json:"is_team_a"`
}
```
