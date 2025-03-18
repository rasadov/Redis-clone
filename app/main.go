package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"time"
)

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}
	fmt.Println("Listening on port 6379...")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		for {
			conn, err := l.Accept()
			if err != nil {
				fmt.Println("Error accepting connection: ", err.Error())
				continue
			}
			go handleConnection(ctx, conn)
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			time.Sleep(100 * time.Second)
		}
	}
}

func handleConnection(ctx context.Context, conn net.Conn) {
	defer conn.Close()
	message := make([]byte, 1024)
	for {
		select {
		case <-ctx.Done():
			return
		default:
			bytesRead, err := conn.Read(message)
			if err != nil {
				fmt.Println("Error reading from connection: ", err.Error())
				return
			}
			if bytesRead != 0 {
				conn.Write([]byte("+PONG\r\n"))
			}
		}
	}
}
