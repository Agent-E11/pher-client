package state

import "github.com/agent-e11/pher-client/internal/menu"

type AppState struct {
	CurrentMenu menu.Menu
	LineNum int
}
