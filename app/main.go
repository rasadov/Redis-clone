package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strings"
)

var (
	Storage *InMemoryStorage
)

func main() {
	Storage = NewInMemoryStorage()
	eventLoop := &EventLoop{
		mainTasks: make(chan Task, 100),
		stop:      make(chan bool),
	}

	wg := InitEventLoop(eventLoop, 5)

	listener, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379:", err)
		return
	}
	defer listener.Close()
	fmt.Println("Listening on port 6379...")

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				fmt.Println("Error accepting connection:", err)
				continue
			}

			// Add connection handling to the event loop
			AddToEventLoop(eventLoop, Task{
				MainTask:   handleConnection,
				IsBlocking: true,
			}, conn)
		}
	}()

	wg.Wait()
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	for {
		array, err := ReadArray(reader)
		if err != nil {
			if err == io.EOF {
				fmt.Println("Client disconnected.")
			} else {
				fmt.Println("Error parsing command:", err)
			}
			return
		}

		if len(array) == 0 {
			continue
		}

		cmd := strings.ToUpper(array[0])
		switch cmd {
		case "PING":
			writeSimpleString(conn, "PONG")

		case "ECHO":
			if len(array) < 2 {
				writeErrorString(conn, "Wrong number of arguments for 'ECHO' command")
				continue
			}
			echoMsg := array[1] // The text to echo
			writeBulkString(conn, echoMsg)

		case "GET":
			key := array[1]
			value, ok := Storage.Get(key)
			if ok {
				writeBulkString(conn, value)
			} else {
				writeSimpleString(conn, "$-1\\r\\n")
			}

		case "SET":
			key, value := array[1], array[2]
			Storage.SetKey(key, value)
			writeBulkString(conn, "OK")

		default:
			writeErrorString(conn, fmt.Sprintf("unknown command '%s'", cmd))
		}
	}
}
