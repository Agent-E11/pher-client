package keybinding

import (
	"slices"

	"github.com/gdamore/tcell/v2"
)

var keyNames = tcell.KeyNames

var keyBindings = map[string][]string{
	"quit": {
		"r" + "q",
		"r" + "Q",
		"k" + keyNames[tcell.KeyCtrlC],
		"k" + keyNames[tcell.KeyEsc],
	},
	"select": {
		"k" + keyNames[tcell.KeyEnter],
	},
	"up": {
		"r" + "k",
		"k" + keyNames[tcell.KeyUp],
	},
	"down": {
		"r" + "j",
		"k" + keyNames[tcell.KeyDown],
	},
}

func IsAction(ek *tcell.EventKey, action string) bool {
	keys, ok := keyBindings[action]
	if !ok {
		return false
	}
	if ek.Key() == tcell.KeyRune {
		return slices.Contains(keys, "r" + string(ek.Rune()))
	}

	return slices.Contains(keys, "k" + keyNames[ek.Key()])
}
