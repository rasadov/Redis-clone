package rdb

import (
	"encoding/binary"
	"github.com/rasadov/redis-clone/app/models"
	"hash/crc64"
	"os"
)

// SaveRDB writes a minimal RDB v11 file to `filename`.
// - Single DB (db=0).
// - Optional expiration in milliseconds (if expiration != zero).
func SaveRDB(filename string, s *models.InMemoryStorage) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	// We'll accumulate the entire RDB in-memory for a final CRC64,
	// or we can do a streaming approach with a rolling CRC.
	// For simplicity, let's just buffer in memory (not great for huge Data).
	var buf buffer

	// 1) Write header: "REDIS" + "0011"
	buf.write([]byte("REDIS0011"))

	// 2) Write metadata: "redis-ver" + "6.0.16"
	// 	  Start metadata section with 0xFA
	buf.writeByte(0xFA)
	encodeString(&buf, "redis-ver")
	encodeString(&buf, "6.0.16")

	// 3) Database section
	//    Start DB opcode: 0xFE
	buf.writeByte(0xFE)

	//    DB index (0). We'll store it as a "size encoded" 0.
	encodeSize(&buf, 0)

	//    Next, 0xFB indicates "hash table size info"
	buf.writeByte(0xFB)

	//    We have N total keys, M keys with expiry
	totalKeys := len(s.Data)
	expires := 0
	for _, e := range s.Data {
		if !e.Expiration.IsZero() {
			expires++
		}
	}
	//    Encode these two sizes
	encodeSize(&buf, uint64(totalKeys))
	encodeSize(&buf, uint64(expires))

	//    Then each key
	for k, e := range s.Data {
		// optional TTL:
		if !e.Expiration.IsZero() {
			// 0xFC => Expire in milliseconds
			buf.writeByte(0xFC)
			// 8-byte little-endian of the Unix time in milliseconds
			ms := e.Expiration.UnixMilli()
			var tmp [8]byte
			binary.LittleEndian.PutUint64(tmp[:], uint64(ms))
			buf.write(tmp[:])
		}

		// value type = 0 => string
		buf.writeByte(0x00)

		// Key (size-encoded string)
		encodeString(&buf, k)

		// Value (size-encoded string)
		encodeString(&buf, e.Value)
	}

	// 4) End of file marker
	buf.writeByte(0xFF)

	// 5) CRC64 of everything so far
	crcTable := crc64.MakeTable(crc64.ISO)
	checksum := crc64.Checksum(buf.bytes, crcTable)
	var crcBuf [8]byte
	binary.LittleEndian.PutUint64(crcBuf[:], checksum)
	buf.write(crcBuf[:])

	// Now flush buf to file
	_, err = f.Write(buf.bytes)
	return err
}
