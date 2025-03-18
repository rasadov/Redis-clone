package main

import (
	"fmt"
	"io"
	"net"
)

// writeSimpleString writes something like: +OK\r\n
func writeString(conn net.Conn, s string) {
	// Just write the raw bytes, no extra escapes/quotes
	conn.Write([]byte(s))
}

func writeSimpleString(w io.Writer, msg string) {
	fmt.Fprintf(w, "+%s\r\n", msg)
}

// writeErrorString writes something like: -ERR <msg>\r\n
func writeErrorString(w io.Writer, msg string) {
	fmt.Fprintf(w, "-ERR %s\r\n", msg)
}

// writeBulkString writes something like: $3\r\nhey\r\n
func writeBulkString(w io.Writer, data string) {
	length := len(data)
	fmt.Fprintf(w, "$%d\r\n%s\r\n", length, data)
}

func writeArray(w io.Writer, data []string) {
	res := fmt.Sprintf("*%d\r\n", len(data))
	for _, val := range data {
		res += fmt.Sprintf("$%d\r\n%s\r\n", len(val), val)
	}
	w.Write([]byte(res))
}
