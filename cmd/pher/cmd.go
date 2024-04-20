package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/agent-e11/pher-client/internal/display"
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

	as := state.AppState{
		LineNum: 0,
	}

	m, err := request.RequestMenu("", address, port)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}

	as.CurrentMenu = m

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

	for {
		display.DrawTextWrap(
			s,
			0, 0 - as.LineNum,
			width - 1, height - 1,
			defStyle,
			as.CurrentMenu.ToString(),
		)
		s.Show()

		ev := s.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEsc ||
				ev.Key() == tcell.KeyCtrlC ||
				ev.Rune() == 'Q' ||
				ev.Rune() == 'q' {
				return
			} else if ev.Key() == tcell.KeyEnter {
				// HACK:
				entity := as.CurrentMenu.DirEntities[as.LineNum]
				m, err := request.RequestMenu(entity.Selector, entity.Hostname, entity.Port)
				if err != nil {
					display.DrawTextWrap(
						s,
						0, 0,
						width - 1, height - 1,
						defStyle,
						fmt.Sprintf("error: %v", err),
					)
					break
				}
				as.CurrentMenu = m
			} else if ev.Rune() == 'j' || ev.Key() == tcell.KeyDown {
				as.LineNum++
			} else if ev.Rune() == 'k' || ev.Key() == tcell.KeyUp {
				as.LineNum--
			}
		}
	}
}
