package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/agent-e11/pher-client/internal/request"
)

func main() {
	var showHelpShort bool
	flag.BoolVar(&showHelpShort, "h", false, "show help")
	var showHelpLong bool
	flag.BoolVar(&showHelpLong, "help", false, "show help")
	var port int
	flag.IntVar(&port, "p", 70, "port number")

	flag.Parse()

	// TODO: Create a `usage` string

	if showHelpShort || showHelpLong {
		fmt.Println("Help menu has not been implemented yet")
		os.Exit(0)
	}

	address := flag.Arg(0)

	if address == "" {
		fmt.Println("please supply an address to connect to")
		os.Exit(1)
	}

	fmt.Println("Address:", address)

	m, err := request.RequestMenu("", address, port)

	m.Debugln()
	fmt.Printf("error: %v\n", err)
}
