package main

import (
	"fmt"
	"net"
	"os"
)

var (
	conn      net.Conn
	l         net.Listener
	err       error
	message   = make([]byte, 1024)
	bytesRead int
)

func main() {

	l, err = net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}
	fmt.Println("Listening on port 6379...")
	conn, err = l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}
	for {
		bytesRead, err = conn.Read(message)
		if err != nil {
			fmt.Println("Error reading from connection: ", err.Error())
			break
		}
		if bytesRead != 0 {
			conn.Write([]byte("+PONG\r\n"))
		}
	}
}
