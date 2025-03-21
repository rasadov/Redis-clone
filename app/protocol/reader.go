package protocol

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strconv"
)

func ReadArray(reader *bufio.Reader) ([]string, error) {
	// Read the first byte, expect '*'
	b, err := reader.ReadByte()
	if err != nil {
		return nil, err
	}
	if b != '*' {
		return nil, fmt.Errorf("expected '*', got '%c'", b)
	}

	line, err := readLine(reader)
	if err != nil {
		return nil, err
	}
	count, err := strconv.ParseInt(line, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid array length: %v", err)
	}

	if count < 0 {
		return nil, nil
	}

	var result []string
	for i := int64(0); i < count; i++ {
		str, err := parseBulkString(reader)
		if err != nil {
			return nil, err
		}
		result = append(result, str)
	}

	return result, nil
}

// parseBulkString expects the next RESP object to be '$<len>\r\n<Data>\r\n'.
func parseBulkString(reader *bufio.Reader) (string, error) {
	b, err := reader.ReadByte()
	if err != nil {
		return "", err
	}
	if b != '$' {
		return "", fmt.Errorf("expected '$', got '%c'", b)
	}

	// Read the length
	line, err := readLine(reader)
	if err != nil {
		return "", err
	}
	length, err := strconv.ParseInt(line, 10, 64)
	if err != nil {
		return "", fmt.Errorf("invalid bulk string length: %v", err)
	}

	if length < 0 {
		return "", nil
	}

	buf := make([]byte, length)
	if _, err := io.ReadFull(reader, buf); err != nil {
		return "", fmt.Errorf("could not read bulk string: %v", err)
	}

	if err := skipCRLF(reader); err != nil {
		return "", err
	}

	return string(buf), nil
}

// readLine reads until \r\n.
func readLine(reader *bufio.Reader) (string, error) {
	raw, err := reader.ReadBytes('\n') // read until newline
	if err != nil {
		return "", err
	}
	// raw should end with '\n', remove it
	raw = bytes.TrimSuffix(raw, []byte("\n"))
	// also remove optional '\r'
	raw = bytes.TrimSuffix(raw, []byte("\r"))
	return string(raw), nil
}

func skipCRLF(reader *bufio.Reader) error {
	b1, err := reader.ReadByte()
	if err != nil {
		return err
	}
	b2, err := reader.ReadByte()
	if err != nil {
		return err
	}
	if b1 != '\r' || b2 != '\n' {
		return fmt.Errorf("expected CRLF, got [%q, %q]", b1, b2)
	}
	return nil
}
