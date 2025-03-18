package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"time"
)

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
			echoMsg := array[1]
			writeBulkString(conn, echoMsg)

		case "GET":
			key := array[1]
			value, ok := Storage.Get(key)
			if ok {
				writeBulkString(conn, value)
			} else {
				writeString(conn, "$-1\r\n")
			}

		case "SET":
			key, value := array[1], array[2]
			fmt.Println("Array ", array)
			withTtl := len(array) > 3 && strings.ToUpper(array[3]) == "PX"
			if withTtl {
				ttl, err := strconv.Atoi(array[4])
				if err != nil {
					writeErrorString(conn, "Error parsing expiration time")
					return
				}
				ttlDuration := time.Duration(ttl) * time.Millisecond
				Storage.SetKeyWithTTL(key, value, ttlDuration)
			} else {
				Storage.SetKey(key, value)
			}
			writeBulkString(conn, "OK")
		case "CONFIG":
			if len(array) < 3 {
				writeErrorString(conn, "Bad request")
				return
			}

			cmd2 := strings.ToUpper(array[2])

			switch cmd2 {
			case "GET":
				if len(array) < 4 {
					return
				}
				key := array[3]
				writeBulkString(conn, Config[key])
				return
			default:
			}

		default:
			writeErrorString(conn, fmt.Sprintf("unknown command '%s'", cmd))
		}
	}
}
