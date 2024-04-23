package main

import (
	"fmt"
	"log"
	"os"

	"github.com/agent-e11/pher-client/internal/display"
	kb "github.com/agent-e11/pher-client/internal/keybinding"
	"github.com/agent-e11/pher-client/internal/request"
	"github.com/agent-e11/pher-client/internal/state"

	"github.com/gdamore/tcell/v2"
	"github.com/spf13/pflag"
)

func main() {
	logFile, err := os.Create("./run.log")
	if err != nil {
		log.Fatalf("error creating logFile: %v", err)
		return
	}
	log.SetOutput(logFile)
	defer logFile.Close()

	// Create and parse flags
	var showHelp bool
	pflag.BoolVarP(&showHelp, "help", "h", false, "show help")
	var port int
	pflag.IntVarP(&port, "port", "p", 70, "port number")
	var selector string
	pflag.StringVarP(&selector, "selector", "s", "", "initial selector string")

	pflag.Parse()

	// Show help and exit
	if showHelp {
		pflag.Usage()
		os.Exit(0)
	}

	address := pflag.Arg(0)

	if address == "" {
		fmt.Println("please supply an address to connect to")
		os.Exit(1)
	}

	fmt.Println("Address:", address)
	fmt.Println("Port:", port)

	state := state.AppState{
		LineNum: -5,
		SelectedIdx: 0,
	}

	m, err := request.RequestMenu(selector, address, port)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}

	state.CurrentMenu = m

	defStyle := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)
	s, err := tcell.NewScreen()
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}
	if err := s.Init(); err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}

	quit := func() {
		maybePanic := recover()
		s.Fini()
		if maybePanic != nil {
			panic(maybePanic)
		}
	}
	defer quit()

	// NOTE: This is for debugging
	msg := ""

	for {
		// Recalculate size every frame
		width, height := s.Size()

		// Number of lines between the cursor and the top/bottom of the screen
		var distToTop, distToBottom =
			display.DisplayMenu(s, state.CurrentMenu, state.SelectedIdx, state.LineNum, 8, nil)

		// Draw debug message
		display.DrawTextWrap(
			s,
			width-len(msg)-1, 0,
			width-1, height-1,
			defStyle,
			fmt.Sprintf(msg),
		)
		msg = ""
		s.Show()

		ev := s.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventKey:
			if kb.IsAction(ev, kb.ActionQuit) {
				return
			} else if kb.IsAction(ev, kb.ActionSelect) {
				// HACK:

				// Get the current selected item
				entity := state.CurrentMenu.DirEntities[state.SelectedIdx]
				// Request a menu using the information in the item
				m, err := request.RequestMenu(entity.Selector, entity.Hostname, entity.Port)
				if err != nil {
					msg += fmt.Sprintf(" error: %v", err)
					break
				}
				// Set the current menu
				state.CurrentMenu = m
				// Reset the screen and cursor position
				state.LineNum = -5
				state.SelectedIdx = 0
			} else if kb.IsAction(ev, kb.ActionDown) && state.LineNum < len(state.CurrentMenu.DirEntities) {
				// Move cursor down one item
				state.SelectedIdx++
				// Move screen down if needed
				if distToBottom <= 5 {
					state.LineNum++
				}
			} else if kb.IsAction(ev, kb.ActionUp) && state.LineNum > -height {
				// Move cursor up one item
				state.SelectedIdx--
				// Move screen up if needed
				if distToTop <= 5 {
					state.LineNum--
				}
			} else if kb.IsAction(ev, kb.ActionCenterScreen) {
				// Center screen on cursor
				state.LineNum = state.SelectedIdx - height/2
			} else if kb.IsAction(ev, kb.ActionUpHalfScreen) {
				// Scroll up by half the screen height
				state.SelectedIdx -= height/2
				if state.SelectedIdx < 0 {
					state.SelectedIdx = 0
				}
				// Center screen on cursor
				state.LineNum = state.SelectedIdx - height/2
			} else if kb.IsAction(ev, kb.ActionDownHalfScreen) {
				// Scroll down by half the screen height
				state.SelectedIdx += height/2
				if state.SelectedIdx >= len(state.CurrentMenu.DirEntities) {
					state.SelectedIdx = len(state.CurrentMenu.DirEntities) - 1
				}
				// Center screen on cursor
				state.LineNum = state.SelectedIdx - height/2
			}
		}
		msg += fmt.Sprintf(" LineNum: %v", state.LineNum)
	}
}
