package main

import (
	"fmt"
	"log"
	"os"

	"github.com/agent-e11/pher-client/internal/display"
	"github.com/agent-e11/pher-client/internal/keybinding"
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

		display.DisplayMenu(s, state.CurrentMenu, state.SelectedIdx, state.LineNum, 8, nil)

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
			if keybinding.IsAction(ev, "quit") {
				return
			} else if keybinding.IsAction(ev, "select") {
				// HACK:
				entity := state.CurrentMenu.DirEntities[state.SelectedIdx]
				m, err := request.RequestMenu(entity.Selector, entity.Hostname, entity.Port)
				if err != nil {
					msg += fmt.Sprintf(" error: %v", err)
					break
				}
				state.CurrentMenu = m
				state.LineNum = -5
				state.SelectedIdx = 0
			} else if keybinding.IsAction(ev, "down") && state.LineNum < len(state.CurrentMenu.DirEntities) {
				state.LineNum++
				state.SelectedIdx++
			} else if keybinding.IsAction(ev, "up") && state.LineNum > -height {
				state.LineNum--
				state.SelectedIdx--
			}
		}
		msg += fmt.Sprintf(" LineNum: %v", state.LineNum)
	}
}
