package protocol

import (
	"fmt"
	"io"
	"net"
)

// WriteString writes same string that was passed
func WriteString(conn net.Conn, s string) {
	conn.Write([]byte(s))
}

// WriteSimpleString writes something like: +OK\r\n
func WriteSimpleString(w io.Writer, msg string) {
	fmt.Fprintf(w, "+%s\r\n", msg)
}

// WriteErrorString writes something like: -ERR <msg>\r\n
func WriteErrorString(w io.Writer, msg string) {
	fmt.Fprintf(w, "-ERR %s\r\n", msg)
}

// WriteBulkString writes something like: $3\r\nhey\r\n
func WriteBulkString(w io.Writer, data string) {
	length := len(data)
	fmt.Fprintf(w, "$%d\r\n%s\r\n", length, data)
}

// WriteArray writes something like: *2\r\n$2\r\nOK\r\n$6\r\nNOT_OK\r\n
func WriteArray(w io.Writer, data []string) {
	res := fmt.Sprintf("*%d\r\n", len(data))
	for _, val := range data {
		res += fmt.Sprintf("$%d\r\n%s\r\n", len(val), val)
	}
	w.Write([]byte(res))
}
