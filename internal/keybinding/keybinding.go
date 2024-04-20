package keybinding

import (
	"slices"

	"github.com/gdamore/tcell/v2"
)

var keytostr = map[tcell.Key]string{
	// TODO: I think these should probably be all uppercase
	tcell.KeyEsc: "<esc>",
	tcell.KeyCtrlC: "<C-c>",
	tcell.KeyEnter: "<cr>",
	tcell.KeyUp: "<up>",
	tcell.KeyDown: "<down>",
}

var keyBindings = map[string][]string{
	"quit": {
		"r" + "q",
		"r" + "Q",
		"k" + keytostr[tcell.KeyCtrlC],
		"k" + keytostr[tcell.KeyEsc],
	},
	"select": {
		"k" + keytostr[tcell.KeyEnter],
	},
	"up": {
		"r" + "k",
		"k" + keytostr[tcell.KeyUp],
	},
	"down": {
		"r" + "j",
		"k" + keytostr[tcell.KeyDown],
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

	return slices.Contains(keys, "k" + keytostr[ek.Key()])
}
