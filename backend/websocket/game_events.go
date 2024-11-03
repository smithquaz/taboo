package websocket

import "taboo-game/types"

// GameEventsHandler defines the interface for handling game events
type GameEventsHandler interface {
	types.GameEventsServiceInterface
}
