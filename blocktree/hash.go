package blocktree

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"log"
)

// Hash bytes
const HASH_BYTES = 256 / 8

// Hash
type Hash [HASH_BYTES]byte

// Nil Hash
var NilHash = Hash{}

// Fill Hash with SHA256 array
func (h *Hash) Read(b [32]byte) {
	copy(h[:], b[:HASH_BYTES])
}

// Int to bytes
func toBytes(num interface{}) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Fatal(err)
	}
	return buff.Bytes()
}

// Encode to string
func (h *Hash) String() string {
	return hex.EncodeToString(h[:])
}
