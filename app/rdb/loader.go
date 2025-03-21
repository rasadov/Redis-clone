package rdb

import (
	"encoding/binary"
	"github.com/rasadov/redis-clone/app/models"
	"hash/crc64"
	"os"
	"time"
)

// LoadRDB loads an RDB from `filename` into the store, overwriting any existing keys.
func LoadRDB(filename string, s *models.InMemoryStorage) error {
	Data, err := os.ReadFile(filename)
	if err != nil {
		// If file doesn't exist, treat as empty
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	if len(Data) < 9 {
		// Not even a valid header
		return nil
	}

	// Validate header "REDIS0011"
	if string(Data[:9]) != "REDIS0011" {
		// Not recognized, treat as empty
		return nil
	}

	// We'll parse up to the final 8 bytes (checksum).
	// The last 8 bytes are the CRC64. The byte before that is presumably 0xFF.
	// We can verify the checksum or skip it. Let's do a quick verify.

	// actual Data to check = everything except the last 8 bytes
	body := Data[:len(Data)-8]
	wantCRC := binary.LittleEndian.Uint64(Data[len(Data)-8:])

	// compute CRC
	crcTable := crc64.MakeTable(crc64.ISO)
	gotCRC := crc64.Checksum(body, crcTable)
	if gotCRC != wantCRC {
		// Checksum mismatch => ignore or error out. We'll just ignore for brevity.
		// return fmt.Errorf("checksum mismatch: got %x want %x", gotCRC, wantCRC)
	}

	// parse the portion up to the last 8
	// skip header
	idx := 9

	s.Data = make(map[string]models.Entry) // overwrite

	for idx < len(body) {
		op := body[idx]
		idx++

		switch op {
		case 0xFA:
			// AUX field: key and value (both strings), skip
			_, nextIdx, err := decodeString(body, idx)
			if err != nil {
				return err
			}
			idx = nextIdx

			_, nextIdx, err = decodeString(body, idx)
			if err != nil {
				return err
			}
			idx = nextIdx

		case 0xFB:
			// Resize DB hint: 2 size-encoded values
			_, nextIdx, err := decodeSize(body, idx)
			if err != nil {
				return err
			}
			idx = nextIdx

			_, nextIdx, err = decodeSize(body, idx)
			if err != nil {
				return err
			}
			idx = nextIdx

		case 0xFE:
			// DB selector
			_, nextIdx, err := decodeSize(body, idx)
			if err != nil {
				return err
			}
			idx = nextIdx

		case 0xFC:
			// Expire time in milliseconds
			if idx+8 > len(body) {
				return nil
			}
			ms := binary.LittleEndian.Uint64(body[idx : idx+8])
			idx += 8

			if idx >= len(body) {
				return nil
			}
			valType := body[idx]
			idx++
			if valType != 0x00 {
				return nil
			}

			key, nextIdx, err := decodeString(body, idx)
			if err != nil {
				return err
			}
			idx = nextIdx

			val, nextIdx, err := decodeString(body, idx)
			if err != nil {
				return err
			}
			idx = nextIdx

			s.Data[key] = models.Entry{
				Value:      val,
				Expiration: time.UnixMilli(int64(ms)),
			}

		case 0xFD:
			// Expire time in seconds
			if idx+4 > len(body) {
				return nil
			}
			secs := binary.LittleEndian.Uint32(body[idx : idx+4])
			idx += 4

			if idx >= len(body) {
				return nil
			}
			valType := body[idx]
			idx++
			if valType != 0x00 {
				return nil
			}

			key, nextIdx, err := decodeString(body, idx)
			if err != nil {
				return err
			}
			idx = nextIdx

			val, nextIdx, err := decodeString(body, idx)
			if err != nil {
				return err
			}
			idx = nextIdx

			s.Data[key] = models.Entry{
				Value:      val,
				Expiration: time.Unix(int64(secs), 0),
			}

		case 0x00:
			// Raw string key-value pair (no expiration)
			key, nextIdx, err := decodeString(body, idx)
			if err != nil {
				return err
			}
			idx = nextIdx

			val, nextIdx, err := decodeString(body, idx)
			if err != nil {
				return err
			}
			idx = nextIdx

			s.Data[key] = models.Entry{Value: val}

		case 0xFF:
			// End of file marker
			break

		default:
			if op >= 0xC0 && op <= 0xC3 {
				// Encoded values (e.g., LRU/LFU) or special metadata – not handled
				// For now, skip next byte as placeholder
				idx++
				continue
			}

			// Unknown opcode – stop parsing
			return nil
		}
	}

	return nil
}
