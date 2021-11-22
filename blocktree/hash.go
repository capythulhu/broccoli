package blocktree

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"log"
)

// Hash bytes
const HASH_BYTES = 32 / 8

// Hash
type Hash [HASH_BYTES]byte

// Nil Hash
var NilHash = Hash{}

// Fill Hash with bytes slice
func DecodeSHA256(b [32]byte) (h Hash) {
	copy(h[:], b[:HASH_BYTES])
	return
}

// Int to bytes
func toBytes(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}
	return buff.Bytes()
}

// Encode to string
func (h *Hash) String() string {
	return hex.EncodeToString(h[:])
}
