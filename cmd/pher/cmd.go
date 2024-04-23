package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/agent-e11/pher-client/internal/display"
	"github.com/agent-e11/pher-client/internal/keybinding"
	"github.com/agent-e11/pher-client/internal/request"
	"github.com/agent-e11/pher-client/internal/state"
	"github.com/gdamore/tcell/v2"
)

// Take short, long, and default flag values, and reconcile conflictions in
// the following way:
//
// If long is not default, use long
//
// If long is default, use short (whether it is default or not)
//
// This gives precedence to the long flag (Because it is more specific).
// Ideally, the precedence would go to the flag that was set furthest to the
// right, but I don't know how I would do that at the moment
func reconcileShortLongFlags[T comparable](short T, long T, defultVal T) T {
	if long != defultVal {
		return long
	}
	return short
}

func main() {
	logFile, err := os.Create("./first.log")
	if err != nil {
		log.Fatalf("error creating logFile: %v", err)
		return
	}
	log.SetOutput(logFile)
	defer logFile.Close()

	var showHelpShort bool
	flag.BoolVar(&showHelpShort, "h", false, "show help")
	var showHelpLong bool
	flag.BoolVar(&showHelpLong, "help", false, "show help")
	var portShort int
	flag.IntVar(&portShort, "p", 70, "port number")
	var portLong int
	flag.IntVar(&portLong, "port", 70, "port number")
	var selectorShort string
	flag.StringVar(&selectorShort, "s", "", "initial selector string")
	var selectorLong string
	flag.StringVar(&selectorLong, "selector", "", "initial selector string")

	flag.Parse()

	// Show help
	if showHelpShort || showHelpLong {
		flag.Usage()
		os.Exit(0)
	}

	port := reconcileShortLongFlags(portShort, portLong, 70)
	selector := reconcileShortLongFlags(selectorShort, selectorLong, "")

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

		display.DisplayMenu(s, state.CurrentMenu, state.LineNum, nil)
		
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
			} else if keybinding.IsAction(ev, "down") && state.LineNum < len(state.CurrentMenu.DirEntities) {
				state.LineNum++
			} else if keybinding.IsAction(ev, "up") && state.LineNum > -height {
				state.LineNum--
			}
		}
		msg += fmt.Sprintf("LineNum: %v", state.LineNum)
	}
}
