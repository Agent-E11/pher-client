package keybinding

import (
	"slices"

	"github.com/gdamore/tcell/v2"
)

var keyNames = tcell.KeyNames

type Action uint

const (
	ActionQuit Action = iota
	ActionSelect
	ActionCenterScreen
	ActionUpHalfScreen
	ActionDownHalfScreen
	ActionOpenHelp
	ActionUp
	ActionDown
	ActionLeft
	ActionRight
)

var keyBindings = map[Action][]string{
	ActionQuit: {
		"r" + "q",
		"r" + "Q",
		"k" + keyNames[tcell.KeyCtrlC],
		"k" + keyNames[tcell.KeyEsc],
	},
	ActionSelect: {
		"k" + keyNames[tcell.KeyEnter],
	},
	ActionCenterScreen: {
		"r" + "z",
	},
	ActionUpHalfScreen: {
		"k" + keyNames[tcell.KeyCtrlU],
	},
	ActionDownHalfScreen: {
		"k" + keyNames[tcell.KeyCtrlD],
	},
	ActionUp: {
		"r" + "k",
		"k" + keyNames[tcell.KeyUp],
	},
	ActionDown: {
		"r" + "j",
		"k" + keyNames[tcell.KeyDown],
	},
}

func IsAction(event *tcell.EventKey, action Action) bool {
	keys, ok := keyBindings[action]
	if !ok {
		return false
	}
	if event.Key() == tcell.KeyRune {
		return slices.Contains(keys, "r" + string(event.Rune()))
	}

	return slices.Contains(keys, "k" + keyNames[event.Key()])
}
