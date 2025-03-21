package main

import (
	"bufio"
	"fmt"
	"github.com/rasadov/redis-clone/app/protocol"
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
		array, err := protocol.ReadArray(reader)
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
			protocol.WriteSimpleString(conn, "PONG")

		case "ECHO":
			if len(array) < 2 {
				protocol.WriteErrorString(conn, "Wrong number of arguments for 'ECHO' command")
				continue
			}
			echoMsg := array[1]
			protocol.WriteBulkString(conn, echoMsg)

		case "GET":
			key := array[1]
			fmt.Println("Key is", key)
			value, ok := Storage.Get(key)
			fmt.Println("Value is", value)
			if ok {
				protocol.WriteBulkString(conn, value)
			} else {
				protocol.WriteString(conn, "$-1\r\n")
			}

		case "SET":
			key, value := array[1], array[2]
			fmt.Println("Array ", array)
			withTtl := len(array) > 3 && strings.ToUpper(array[3]) == "PX"
			if withTtl {
				ttl, err := strconv.Atoi(array[4])
				if err != nil {
					protocol.WriteErrorString(conn, "Error parsing expiration time")
					return
				}
				ttlDuration := time.Duration(ttl) * time.Millisecond
				Storage.SetKeyWithTTL(key, value, ttlDuration)
			} else {
				Storage.SetKey(key, value)
			}
			protocol.WriteBulkString(conn, "OK")
		case "CONFIG":
			cmd2 := strings.ToUpper(array[1])

			switch cmd2 {
			case "GET":
				key := array[2]
				value, _ := Config[key]
				protocol.WriteArray(conn, []string{key, value})
				return
			default:
			}
		case "KEYS":
			pattern := array[1]
			keys := Storage.Keys(pattern)
			protocol.WriteArray(conn, keys)
		default:
			protocol.WriteErrorString(conn, fmt.Sprintf("unknown command '%s'", cmd))
		}
	}
}
