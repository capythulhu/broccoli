package blocktree

import (
	"bytes"
	"encoding/binary"
	"log"
)

// Interface to bytes
func toBytes(num interface{}) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Fatal(err)
	}
	return buff.Bytes()
}
