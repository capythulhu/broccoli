package blocktree

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"log"
)

// Hash
type Hash [32]byte

// Nil Hash
var NilHash = Hash{}

// Fill Hash with bytes slice
func DecodeHash(b []byte) (h Hash) {
	copy(h[:], b[:32])
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
