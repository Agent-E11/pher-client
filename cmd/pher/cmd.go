package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/agent-e11/pher-client/menu"
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

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", address, port))
	if err != nil {
		fmt.Printf("error: %v", err)
		os.Exit(1)
	}
	defer conn.Close()

	buf := make([]byte, 1024)

	// Write an empty line, indicating to the server to "list what you have"
	// RFC 1436 Page 3, Paragraph 1 (?)
	conn.Write([]byte("\r\n"))

	menuString := ""

	for {
		// Set a timeout for 5 seconds from now
		conn.SetReadDeadline(time.Now().Add(5 * time.Second))
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Printf("error: %v\n", err)
			break
		}

		menuString += string(buf[:n])

		fmt.Printf("Received: %s", buf[:n])
	}

	fmt.Println("Done")

	fmt.Print("\n\n\n\n\n\n\n")

	m, err := menu.FromString(menuString)

	m.Debugln()
	fmt.Printf("error: %v\n", err)
}
