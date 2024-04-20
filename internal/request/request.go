package request

import (
	"errors"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/agent-e11/pher-client/internal/menu"
)

var ErrConnect error = errors.New("error connecting to host")
var ErrTimeout error = errors.New("connection timed out")

const rwTimeOut = 5 * time.Second

func RequestMenu(selector string, address string, port int) (m menu.Menu, err error) {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", address, port))
	if err != nil {
		return m, ErrConnect
	}
	defer conn.Close()

	conn.Write([]byte(fmt.Sprintf("%s\r\n", selector)))

	menuString := ""
	buf := make([]byte, 1024)

	for {
		conn.SetReadDeadline(time.Now().Add(rwTimeOut))

		n, err := conn.Read(buf)
		if err == io.EOF {
			break
		// Check if the error is a timeout error
		} else if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			return m, ErrTimeout
		} else if err != nil {
			fmt.Printf("error: %v\n", err)
			return m, err
		}

		menuString += string(buf[:n])
	}

	return menu.FromString(menuString)
}
