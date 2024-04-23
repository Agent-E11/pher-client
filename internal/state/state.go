package state

import "github.com/agent-e11/pher-client/internal/menu"

type AppState struct {
	CurrentMenu menu.Menu
	LineNum int
	// Index of selected item in menu
	SelectedIdx int
}
