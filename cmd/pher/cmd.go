package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/agent-e11/pher-client/internal/display"
	"github.com/agent-e11/pher-client/internal/keybinding"
	"github.com/agent-e11/pher-client/internal/request"
	"github.com/agent-e11/pher-client/internal/state"
	"github.com/gdamore/tcell/v2"
)

func main() {
	var showHelpShort bool
	flag.BoolVar(&showHelpShort, "h", false, "show help")
	var showHelpLong bool
	flag.BoolVar(&showHelpLong, "help", false, "show help")
	var portShort int
	flag.IntVar(&portShort, "p", 70, "port number")
	var portLong int
	flag.IntVar(&portLong, "port", 70, "port number")
	// TODO: Add a --selector flag

	flag.Parse()

	// Show help
	if showHelpShort || showHelpLong {
		flag.Usage()
		os.Exit(0)
	}

	// Set port to portLong if it is not the default
	// else, set it to portShort
	var port int
	if portLong != 70 {
		port = portLong
	} else {
		port = portShort
	}

	address := flag.Arg(0)

	if address == "" {
		fmt.Println("please supply an address to connect to")
		os.Exit(1)
	}

	fmt.Println("Address:", address)
	fmt.Println("Port:", port)

	state := state.AppState{
		LineNum: 0,
	}

	m, err := request.RequestMenu("", address, port)
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

	width, height := s.Size()

	// NOTE: This is for debugging
	msg := ""

	for {
		display.DrawTextWrap(
			s,
			0, 0 - state.LineNum,
			width - 1, height - 1,
			defStyle,
			state.CurrentMenu.ToString(),
		)
		display.DrawTextWrap(
			s,
			0, 0,
			width - 1, height - 1,
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
				entity := state.CurrentMenu.DirEntities[state.LineNum]
				m, err := request.RequestMenu(entity.Selector, entity.Hostname, entity.Port)
				if err != nil {
					msg = fmt.Sprintf("error: %v", err)
					break
				}
				state.CurrentMenu = m
				state.LineNum = 0
			} else if keybinding.IsAction(ev, "down") {
				state.LineNum++
			} else if keybinding.IsAction(ev, "up") {
				state.LineNum--
			}
		}
	}
}
